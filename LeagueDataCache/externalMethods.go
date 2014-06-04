package LeagueDataCache

import (
	"fmt"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

// Public Functions that use channel to communicate with cache goroutine

// Fetch Champion from cache goroutine.
func GetChampion(id int64, cacheRequests chan Request) (c lapi.Champion, e error) {
	result := make(chan Response, 1)
	cReq := Request{Type: "champion", Key: fmt.Sprintf("%d", id), Response: result}
	cacheRequests <- cReq
	champResponse := <-result
	if champResponse.Ok {
		champ, _ := champResponse.Value.(lapi.Champion)
		return champ, nil
	}

	return lapi.Champion{}, CacheError{message: "Failed to retrieve champion."}
}

// Fetch Summoner from cache goroutine
func GetSummoner(name string, cache chan Request) (lapi.Summoner, error) {
	value, getErr := goGet(Request{Type: "summoner", Key: name}, cache)
	if getErr != nil {
		return lapi.Summoner{}, CacheError{message: "Failed to retrieve summoner."}
	}
	s, _ := value.(lapi.Summoner)
	return s, nil
}

// Fetch Summoner recent match history typed version
func GetSummonerMatchesSimple(id int64, cache chan Request) (MatchHistory, error) {
	value, getErr := goGet(Request{Type: "games", Key: fmt.Sprintf("%d", id)}, cache)
	if getErr != nil {
		return MatchHistory{}, CacheError{message: "Failed to retrieve games."}
	}
	lmh, _ := value.(MatchHistory)
	return lmh, nil
}

// Typed cache get for a match
func GetMatch(id int64, cache chan Request) (MatchDetail, error) {
	value, getErr := goGet(Request{Type: "game", Key: fmt.Sprintf("%d", id)}, cache)

	if getErr != nil {
		return MatchDetail{}, CacheError{message: "Failed to retrieve game."}
	}
	game, _ := value.(MatchDetail)
	return game, nil
}

// Generic function to request data from cache.
func goGet(request Request, cache chan Request) (interface{}, error) {
	cResponse := make(chan Response, 1)
	request.Response = cResponse
	cache <- request
	respValue := <-cResponse
	if respValue.Ok {
		return respValue.Value, nil
	}
	return nil, CacheError{message: "Failed to fetch from cache!"}
}
