package LeagueDataCache

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
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
	key := datastore.NewKey(mp.Context, "Summoner", "", s.Id, nil)
	_, err := datastore.Put(mp.Context, key, &s)
	if err != nil {
		return err
	}
	return nil
}

func (mp *MemcachePersistance) GetSummoner(s *lapi.Summoner) error {
	return datastore.Get(mp.Context, datastore.NewKey(mp.Context, "Summoner", "", s.Id, nil), &s)
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
