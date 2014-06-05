package LeagueDataCache

import (
	"appengine/aetest"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"testing"
)

func TestPutCache(t *testing.T) {
	// Setup Test
	allSummonersById = make(map[int64]lapi.Summoner)
	allSummonersByName = make(map[string]lapi.Summoner)
	// Build Data
	testResp := Response{Type: "summoner", Value: lapi.Summoner{Name: "Test Summoner", Id: int64(1)}}
	putCache(testResp)
	// Assert
	if _, ok := allSummonersById[int64(1)]; !ok {
		t.Log("Failed to get summoner by Id.")
		t.FailNow()
	}
	if _, ok := allSummonersByName["testsummoner"]; !ok {
		t.Log("Failed to get summoner by Id.")
		t.FailNow()
	}
}

func TestFetchCache(t *testing.T) {
	testSummoner := lapi.Summoner{Name: "Test Summoner", Id: int64(1)}
	// Setup Test
	allSummonersById = map[int64]lapi.Summoner{int64(1): testSummoner}
	allSummonersByName = map[string]lapi.Summoner{"testsummoner": testSummoner}
	responseChannel := make(chan Response, 1)
	// Build Data
	context, _ := aetest.NewContext(nil)
	testReq := Request{Type: "summoner", Key: "Test Summoner", Response: responseChannel, Context: context}
	fetchCache(testReq)
	// Assert
	cacheResponse := <-responseChannel
	if !cacheResponse.Ok {
		t.Log("FetchCache failed to retrieve summoner.")
		t.FailNow()
	}
	if cachedSummoner, ok := cacheResponse.Value.(lapi.Summoner); ok {
		if cachedSummoner.Name != testSummoner.Name {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}
