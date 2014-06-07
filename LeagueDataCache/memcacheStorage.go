package LeagueDataCache

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"errors"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
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
	_, err := datastore.Put(mp.Context, key, &s)
	if err != nil {
		mp.Context.Infof("Failed Storing Summoner: %s", err.Error())
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
	entities := make([]lapi.Summoner, len(keys))
	err := datastore.GetMulti(mp.Context, keys, entities)
	return entities, err
}

func (mp *MemcachePersistance) GetSummonerByName(s *lapi.Summoner) error {
	query := datastore.NewQuery("Summoner").Filter("Name =", s.Name).Limit(1)
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
	return nil
}

func (mp *MemcachePersistance) PutObject(objType string, id string, thing interface{}) error {
	jsonData, marshErr := json.Marshal(thing)
	if marshErr != nil {
		return marshErr
	}
	key := datastore.NewKey(mp.Context, objType, id, 0, nil)
	_, err := datastore.Put(mp.Context, key, &cachedObject{Data: jsonData})
	if err != nil {
		mp.Context.Warningf("Failed to store object: %s", err.Error())
		return err
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
