package LeagueDataCache

import (
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
