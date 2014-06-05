package LeagueDataCache

import (
	"appengine/aetest"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"testing"
)

func TestPutSummonerCache(t *testing.T) {
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

func TestFetchSummonerCache(t *testing.T) {
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

func TestPutGamesCache(t *testing.T) {
	// Setup Test
	allGames = make(map[MatchKey]lapi.Game)
	gamesBySummoner = make(map[int64][]lapi.Game)
	// Build Data
	gameList := []lapi.Game{lapi.Game{ChampionId: int64(10), GameId: int64(3), Stats: lapi.RawStats{}}}
	games := lapi.RecentGames{SummonerId: int64(1), Games: gameList}
	testResp := Response{Type: "games", Value: games}
	putCache(testResp)
	// Assert
	if value, ok := allGames[MatchKey{SummonerId: int64(1), MatchId: int64(3)}]; ok {
		if value.ChampionId != int64(10) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}

//func TestFetchGamesCache(t *testing.T) {
//	testSummoner := lapi.Summoner{Name: "Test Summoner", Id: int64(1)}
//	// Setup Test
//	allSummonersById = map[int64]lapi.Summoner{int64(1): testSummoner}
//	allSummonersByName = map[string]lapi.Summoner{"testsummoner": testSummoner}
//	responseChannel := make(chan Response, 1)
//	// Build Data
//	context, _ := aetest.NewContext(nil)
//	testReq := Request{Type: "summoner", Key: "Test Summoner", Response: responseChannel, Context: context}
//	fetchCache(testReq)
//	// Assert
//	cacheResponse := <-responseChannel
//	if !cacheResponse.Ok {
//		t.Log("FetchCache failed to retrieve summoner.")
//		t.FailNow()
//	}
//	if cachedSummoner, ok := cacheResponse.Value.(lapi.Summoner); ok {
//		if cachedSummoner.Name != testSummoner.Name {
//			t.FailNow()
//		}
//	} else {
//		t.FailNow()
//	}
//}
