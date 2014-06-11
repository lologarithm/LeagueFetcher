package LeagueDataCache

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"errors"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"time"
)

type MemcachePersistance struct {
	Context appengine.Context
}

type cachedObject struct {
	Data       []byte
	IntIndex   int64 // Matches this is the summoner id, Summoners this is id
	CachedDate int64 // Lets you check within a certain date.
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

func (mp *MemcachePersistance) PutObject(objType string, keyName string, index int64, thing interface{}) error {
	jsonData, marshErr := json.Marshal(thing)
	if marshErr != nil {
		return marshErr
	}
	key := datastore.NewKey(mp.Context, objType, keyName, 0, nil)
	_, err := datastore.Put(mp.Context, key, &cachedObject{Data: jsonData, IntIndex: index, CachedDate: getExpireTime(true)})
	if err != nil {
		mp.Context.Warningf("Failed to store object: %s", err.Error())
		return err
	}
	return nil
}

func (mp *MemcachePersistance) PutObjects(objType string, keys []string, indexes []int64, things []interface{}) error {
	cacheThings := make([]*cachedObject, len(things))
	dsKeys := make([]*datastore.Key, len(keys))
	for ind, t := range things {
		jsonData, marshErr := json.Marshal(t)
		if marshErr != nil {
			return marshErr
		}
		cacheThings[ind] = &cachedObject{Data: jsonData, IntIndex: indexes[ind], CachedDate: getExpireTime(true)}
		dsKeys[ind] = datastore.NewKey(mp.Context, objType, keys[ind], 0, nil)
	}

	_, dsErr := datastore.PutMulti(mp.Context, dsKeys, cacheThings)
	if dsErr != nil {
		mp.Context.Warningf("Failed to store multi objects %s: %s", objType)
	}
	return nil
}

func (mp *MemcachePersistance) GetObject(objType string, id string, thing interface{}) error {
	var co cachedObject
	getErr := datastore.Get(mp.Context, datastore.NewKey(mp.Context, objType, id, 0, nil), &co)
	if getErr != nil {
		mp.Context.Warningf("Get Failed: %s", getErr.Error())
		return getErr
	}
	mErr := json.Unmarshal(co.Data, thing)
	if mErr != nil {
		mp.Context.Warningf("Stored Json Failed: %s", mErr.Error())
		return mErr
	}
	return nil
}

func (mp *MemcachePersistance) GetMatchesByIndex(index int64) ([]lapi.Game, error) {
	games := []lapi.Game{}
	query := datastore.NewQuery("Match").Filter("IntIndex =", index).Order("CachedDate").Limit(10)
	var cachedGames []cachedObject
	_, err := query.GetAll(mp.Context, &games)
	if err != nil {
		return nil, err
	} else if len(cachedGames) == 0 {
		return nil, errors.New("No games found.")
	}
	for _, co := range cachedGames {
		if co.CachedDate < time.Now().UnixNano() {
			continue
		}
		var game lapi.Game
		mErr := json.Unmarshal(co.Data, &game)
		if mErr != nil {
			mp.Context.Warningf("Stored Json Failed: %s", mErr.Error())
			return nil, mErr
		}
		games = append(games, game)
	}
	return games, nil
}

//func (mp *MemcachePersistance) GetObjectsByIndex(objType string, index int64, thing interface{}) error {
//	query := datastore.NewQuery("Match").Filter("IntIndex =", index).Order("CachedDate").Limit(10)
//	var co cachedObject
//	getErr := datastore.Get(mp.Context, datastore.NewKey(mp.Context, objType, id, 0, nil), &co)
//	if getErr != nil {
//		mp.Context.Warningf("Get Failed: %s", getErr.Error())
//		return getErr
//	}
//	mErr := json.Unmarshal(co.Data, thing)
//	if mErr != nil {
//		mp.Context.Warningf("Stored Json Failed: %s", mErr.Error())
//		return mErr
//	}
//	return nil
//}
