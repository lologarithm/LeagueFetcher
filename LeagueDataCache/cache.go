// Caching layer for use with the LeagueApi package.
package LeagueDataCache

import (
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"strconv"
)

const (
	championCache = "ccache.json"
	summonerCache = "scache.json"
)

type Request struct {
	Response chan Response
	Type     string
	Key      string
}

type Response struct {
	Value interface{}
	Ok    bool
}

var allSummonersByName map[string]lapi.Summoner
var allSummonersById map[int64]lapi.Summoner
var allChampions map[int64]lapi.Champion
var allLeagues map[int64]lapi.League
var allTeams map[int64]lapi.Team
var allGames map[int64]lapi.Game
var gamesBySummoner map[int64][]lapi.Game

// RunCache is the primary method. Start this as a goroutine and then use other public methods to fetch from this.
func RunCache(exit chan bool, requests chan Request) {
	setupCache()
	for {
		select {
		case <-exit:
			return
		case cRequest := <-requests:
			fetchCache(cRequest)
		}
	}
	saveCache()
}

func fetchCache(request Request) {
	response := Response{Ok: false}
	switch request.Type {
	case "summoner":
		summoner := fetchAndCacheSummoner(request.Key)
		response.Ok = true
		response.Value = summoner
	case "champion":
		intKey, err := strconv.ParseInt(request.Key, 10, 64)
		if err != nil {
			break
		}
		champ := fetchAndCacheChampion(intKey)
		response.Value = champ
		response.Ok = true
	case "games":
		intKey, err := strconv.ParseInt(request.Key, 10, 64)
		if err != nil {
			break
		}
		games := fetchSimpleMatchHistory(intKey)
		response.Value = games
		response.Ok = true
	case "game":
		intKey, err := strconv.ParseInt(request.Key, 10, 64)
		if err != nil {
			break
		}
		game := fetchMatch(intKey)
		response.Value = game
		response.Ok = true
	}

	request.Response <- response
}

func setupCache() {
	loadChampions(championCache)
	loadSummoners(summonerCache)
	allGames = make(map[int64]lapi.Game, 1)
	gamesBySummoner = make(map[int64][]lapi.Game, 1)
}

func saveCache() {
	storeChampions(championCache)
	storeSummoners(summonerCache)
}

type CacheError struct {
	message string
}

func (c CacheError) Error() string {
	return c.message
}
