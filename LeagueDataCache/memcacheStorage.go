package LeagueDataCache

import (
	"encoding/json"
	"errors"
	"time"

	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"

	"appengine"
	"appengine/datastore"
)

type MemcachePersistance struct {
	Context appengine.Context
}

func (mp *MemcachePersistance) PutSummoner(s lapi.Summoner) error {
	key := datastore.NewKey(mp.Context, "TSummoner", "", s.Id, nil)
	s.NormalizedName = NormalizeString(s.Name)
	_, err := datastore.Put(mp.Context, key, &s)
	if err != nil {
		return err
	}
	return nil
}

func (mp *MemcachePersistance) GetSummoner(s *lapi.Summoner) error {
	key := datastore.NewKey(mp.Context, "TSummoner", "", s.Id, nil)
	ferr := datastore.Get(mp.Context, key, s)
	if ferr != nil {
		return ferr
	}
	return nil
}

func (mp *MemcachePersistance) GetSummoners(ids []int64) ([]lapi.Summoner, error) {
	keys := make([]*datastore.Key, len(ids))
	for index, id := range ids {
		keys[index] = datastore.NewKey(mp.Context, "TSummoner", "", id, nil)
	}
	var entities = make([]lapi.Summoner, len(keys))
	err := datastore.GetMulti(mp.Context, keys, entities)

	if err != nil {
		actualEntities := []lapi.Summoner{}
		if me, ok := err.(appengine.MultiError); ok {
			for i, merr := range me {
				if merr == nil {
					actualEntities = append(actualEntities, entities[i])
				}
			}
		}
		return actualEntities, err
	}
	return entities, err
}

func (mp *MemcachePersistance) GetSummonerByName(s *lapi.Summoner) error {
	query := datastore.NewQuery("TSummoner").Filter("NormalizedName =", s.NormalizedName).Limit(1)
	var summoners []lapi.Summoner
	_, err := query.GetAll(mp.Context, &summoners)
	if err != nil {
		return err
	} else if len(summoners) == 0 {
		return errors.New("No summoner found.")
	}
	mp.Context.Infof("FOUND SOMETHING")
	s.Id = summoners[0].Id
	s.ProfileIconId = summoners[0].ProfileIconId
	s.RevisionDate = summoners[0].RevisionDate
	s.SummonerLevel = summoners[0].SummonerLevel
	s.Name = summoners[0].Name
	return nil
}

func (mp *MemcachePersistance) PutMatchDetail(mKey MatchKey, md MatchDetail) error {
	cacheMatch := md.toCachedMatch(mKey.SummonerId)
	dKey := datastore.NewKey(mp.Context, "TMatch", mKey.String(), 0, nil)
	_, err := datastore.Put(mp.Context, dKey, &cacheMatch)
	if err != nil {
		return err
	}
	return nil
}

func (mp *MemcachePersistance) GetMatchDetail(mKey MatchKey, md *MatchDetail) error {
	dKey := datastore.NewKey(mp.Context, "TMatch", mKey.String(), 0, nil)
	var cm cachedMatchDetail
	err := datastore.Get(mp.Context, dKey, &cm)
	if err != nil {
		return err
	}
	e := json.Unmarshal(cm.Data, md)
	if e != nil {
		mp.Context.Warningf("Failed to unmarshal JSON of MatchDetail: %s.\n", e.Error())
	}
	return e
}

func (mp *MemcachePersistance) PutMatchDetails(summonerId int64, matches []MatchDetail) error {
	cachedMatches := make([]cachedMatchDetail, len(matches))
	cacheKeys := make([]*datastore.Key, len(matches))
	for ind, match := range matches {
		cachedMatches[ind] = match.toCachedMatch(summonerId)
		cacheKeys[ind] = datastore.NewKey(mp.Context, "TMatch", cachedMatches[ind].KeyString(), 0, nil)
	}

	_, err := datastore.PutMulti(mp.Context, cacheKeys, cachedMatches)
	if err != nil {
		return err
	}
	return nil
}

func (mp *MemcachePersistance) getCachedMatches(id int64) ([]cachedMatchDetail, error) {
	query := datastore.NewQuery("TMatch").Filter("SummonerId =", id).Order("-PlayedDate").Limit(10)
	var cachedGames []cachedMatchDetail
	_, err := query.GetAll(mp.Context, &cachedGames)
	if err != nil {
		return nil, err
	} else if len(cachedGames) == 0 {
		return nil, errors.New("No games found.")
	} else if cachedGames[0].CacheExpireDate < time.Now().Unix() {
		return cachedGames, errors.New("Games cache expired.")
	}

	return cachedGames, nil
}

func (mp *MemcachePersistance) GetMatchDetails(id int64) ([]MatchDetail, error) {
	games := []MatchDetail{}
	cachedGames, gErr := mp.getCachedMatches(id)
	if gErr != nil {
		return games, gErr
	}
	for _, co := range cachedGames {
		md, cErr := co.ToMatchDetail()
		if cErr != nil {
			mp.Context.Warningf("Stored Json Failed: %s", cErr.Error())
			return nil, cErr
		}
		games = append(games, md)
	}
	return games, nil
}

func (mp *MemcachePersistance) GetMatchHistory(id int64) (mh MatchHistory, e error) {
	cachedGames, e := mp.getCachedMatches(id)
	if e != nil {
		return
	}
	mh, e = convertCachedMatchDetailsToHistory(cachedGames)
	return
}
