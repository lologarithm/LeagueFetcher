package LeagueDataCache

import (
	"appengine/aetest"
	"errors"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"testing"
)

func TestPutSummonerCache(t *testing.T) {
	// Setup Test
	allSummonersById = make(map[int64]lapi.Summoner)
	allSummonersByName = make(map[string]lapi.Summoner)
	// Build Data
	testResp := Response{Type: "summoner", Value: lapi.Summoner{Name: "Test Summoner", Id: int64(1)}, Persist: &MockPersist{}}
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
	testReq := Request{Type: "summoner", Key: "Test Summoner", Response: responseChannel, Context: context, Persist: &MockPersist{}}
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
	gameList := []lapi.Game{lapi.Game{ChampionId: int64(10), GameId: int64(3), Stats: lapi.RawStats{}}, lapi.Game{ChampionId: int64(11), GameId: int64(4), Stats: lapi.RawStats{}}}
	games := lapi.RecentGames{SummonerId: int64(1), Games: gameList}
	testResp := Response{Type: "games", Value: games, Persist: &MockPersist{}}
	putCache(testResp)
	// Assert
	if value, ok := allGames[MatchKey{SummonerId: int64(1), MatchId: int64(3)}]; ok {
		if value.ChampionId != int64(10) {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
	if fetchedGames, ok := gamesBySummoner[int64(1)]; ok {
		if len(fetchedGames) != 2 {
			t.Log("Wrong number of games!")
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}

func TestExternalGetMatch(t *testing.T) {
	SetupCache()
	allChampions = map[int64]lapi.Champion{int64(10): lapi.Champion{Name: "TestChamp1"}, int64(11): lapi.Champion{Name: "TestChamp2"}}
	// Build Data
	gameList := []lapi.Game{lapi.Game{ChampionId: int64(10), GameId: int64(3), Stats: lapi.RawStats{}}, lapi.Game{ChampionId: int64(11), GameId: int64(4), Stats: lapi.RawStats{}}}
	games := lapi.RecentGames{SummonerId: int64(1), Games: gameList}
	testResp := Response{Type: "games", Value: games, Persist: &MockPersist{}}
	putCache(testResp)
	exit := make(chan bool, 1)
	get := make(chan Request, 10)
	put := make(chan Response, 10)
	go RunCache(exit, get, put)
	context, _ := aetest.NewContext(nil)
	match, matchErr := GetMatch(int64(4), int64(1), get, put, context, &MockPersist{})
	if matchErr != nil {
		t.Errorf("Failed to get match: %s", matchErr.Error())
		t.FailNow()
	}
	if match.ChampionName != "TestChamp2" {
		t.Errorf("Wrong Champion Name")
		t.FailNow()
	}

	// Now fetch a match that isn't there.
	_, noMatchErr := GetMatch(int64(6), int64(1), get, put, context, &MockPersist{})
	if noMatchErr == nil {
		t.Log("Found match that doesn't exist.")
		t.FailNow()
	}
}

type MockPersist struct {
}

func (mp *MockPersist) PutObject(objType string, id string, thing interface{}) error {
	return errors.New("Failed to get!")
}

func (mp *MockPersist) GetObject(objType string, id string, thing interface{}) error {
	return errors.New("Failed to get!")
}
