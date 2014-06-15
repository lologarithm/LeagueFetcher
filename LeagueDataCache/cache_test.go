package LeagueDataCache

import (
	"errors"
	"testing"

	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"

	"appengine/aetest"
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
	testReq := Request{Type: "summoner", Key: "Test Summoner", Response: responseChannel}
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
	allGames = make(map[MatchKey]cachedMatchDetail)
	gamesBySummoner = make(map[int64][]cachedMatchDetail)
	// Build Data
	gameList := []lapi.Game{lapi.Game{ChampionId: int64(10), GameId: int64(3), Stats: lapi.RawStats{}}, lapi.Game{ChampionId: int64(11), GameId: int64(4), Stats: lapi.RawStats{}}}
	games := lapi.RecentGames{SummonerId: int64(1), Games: gameList}
	testResp := Response{Type: "games", Value: games}
	putCache(testResp)
	// Assert
	if cg, ok := allGames[MatchKey{SummonerId: int64(1), MatchId: int64(3)}]; ok {
		if cg.GameId != int64(3) {
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
	testResp := Response{Type: "games", Value: games}
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

func TestExternalGetSimpleMatches(t *testing.T) {
	SetupCache()
	allChampions = map[int64]lapi.Champion{int64(10): lapi.Champion{Name: "TestChamp1"}, int64(11): lapi.Champion{Name: "TestChamp2"}}
	exit := make(chan bool, 1)
	get := make(chan Request, 10)
	put := make(chan Response, 10)
	go RunCache(exit, get, put)
	context, _ := aetest.NewContext(nil)
	matches, matchErr := GetSummonerMatchesSimple(int64(1), get, put, context, &MockPersist{})
	if matchErr != nil {
		t.Errorf("Failed to get match: %s", matchErr.Error())
		t.FailNow()
	}
	if len(matches.Games) == 0 {
		t.Errorf("No Matches?: %s", matchErr.Error())
		t.FailNow()
	}
	testKey := MatchKey{MatchId: matches.Games[0].GameId, SummonerId: matches.SummonerId}
	if value, ok := allGames[testKey]; !ok {
		t.Errorf("Couldn't find match: %s", testKey)
		t.FailNow()
	} else {
		t.Logf("Found Match: %v\n", string(value.Data))
	}
}

func TestCacheExpire(t *testing.T) {
	//1. Setup 2 games with expire dates in the past
	//2. Make sure that it counts as expired.

	// 3. Setup 11 games and make sure the last one expires in the future.
	// 4. Make sure it counts as not expired

	// 5. Setup 11 games with last one expires in the past
	// 6. Make sure it counts as expired.
}

type MockPersist struct {
}

func (mp *MockPersist) PutSummoner(s lapi.Summoner) error {
	return errors.New("Failed!")
}
func (mp *MockPersist) GetSummoner(s *lapi.Summoner) error {
	return errors.New("Failed!")
}
func (mp *MockPersist) GetSummonerByName(s *lapi.Summoner) error {
	return errors.New("Failed!")
}
func (mp *MockPersist) GetSummoners(s []int64) ([]lapi.Summoner, error) {
	return nil, errors.New("Failed!")
}
func (mp *MockPersist) GetMatchDetail(a MatchKey, b *MatchDetail) error {
	return errors.New("Failed!")
}
func (mp *MockPersist) PutMatchDetail(a MatchKey, b MatchDetail) error {
	return errors.New("Failed!")
}

func (mp *MockPersist) PutMatchDetails(a int64, b []MatchDetail) error {
	return errors.New("Failed!")
}
func (mp *MockPersist) GetMatchHistory(a int64) (MatchHistory, error) {
	return MatchHistory{}, errors.New("Failed!")
}

func (mp *MockPersist) GetMatchDetails(int64) ([]MatchDetail, error) {
	return []MatchDetail{}, nil
}
