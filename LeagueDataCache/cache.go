// Caching layer for use with the LeagueApi package.
package LeagueDataCache

import (
	"fmt"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"strconv"
	"time"
)

const (
	championCache       = "ccache.json"
	summonerCache       = "scache.json"
	cacheTimeoutMinutes = 30
)

type Request struct {
	Response chan Response
	Type     string
	Key      string
}

type Response struct {
	Value interface{}
	Type  string
	Ok    bool
}

// Cache dictionaries.
var allSummonersByName map[string]lapi.Summoner
var allSummonersById map[int64]lapi.Summoner
var allChampions map[int64]lapi.Champion
var allLeagues map[int64]lapi.League
var allTeams map[int64]lapi.Team
var allGames map[int64]lapi.Game
var gamesBySummoner map[int64][]lapi.Game
var allRankedData map[int64]SummonerRankedData

// RunCache is the primary method. Start this as a goroutine and then use other public methods to fetch from this.
func RunCache(exit chan bool, get chan Request, put chan Response) {
	setupCache()
	for {
		select {
		case <-exit:
			return
		case getObject := <-get:
			fetchCache(getObject)
		case putObject := <-put:
			putCache(putObject)
		}
	}
	saveCache()
}

func putCache(resp Response) {
	switch resp.Type {
	case "summoner":
		if summoner, ok := resp.Value.(lapi.Summoner); ok {
			allSummonersById[summoner.Id] = summoner
			allSummonersByName[summoner.Name] = summoner
		}
	case "champion":
		// For now the cache will handle this.
	case "games":
		if matches, ok := resp.Value.(lapi.RecentGames); ok {
			for _, match := range matches.Games {
				match.ExpireTime = getExpireTime()
				// Don't put the match into the list of games for the summoner if it's already cached.
				if _, ok := allGames[match.GameId]; !ok {
					gamesBySummoner[matches.SummonerId] = append([]lapi.Game{match}, gamesBySummoner[matches.SummonerId]...)
				}
				// Always re-cache here for updated match time.
				allGames[match.GameId] = match
			}
		}
	case "game":
		if game, ok := resp.Value.(lapi.Game); ok {
			allGames[game.GameId] = game
		}
	case "rankedData":
		if data, ok := resp.Value.(SummonerRankedData); ok {
			data.ExpireTime = getExpireTime()
			allRankedData[data.Id] = data
		}
	}
}

func fetchCache(request Request) {
	response := Response{Ok: false}
	switch request.Type {
	case "summoner":
		if summoner, ok := allSummonersByName[request.Key]; ok {
			response.Ok = true
			response.Value = summoner
		}
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
		// If last fetch of game history is old, refetch game history.
		if games, ok := gamesBySummoner[intKey]; ok {
			if len(games) > 0 && games[0].ExpireTime > time.Now().Unix() {
				sliceEnd := 10
				if len(games) < 10 {
					sliceEnd = len(games)
				}
				matchHistory := convertGamesToMatchHistory(intKey, games[0:sliceEnd], fetchAndCacheChampion)
				response.Value = matchHistory
				response.Ok = true
			} else {
				fmt.Printf("Cached games are old or no games found.")
			}
		}
	case "game":
		intKey, err := strconv.ParseInt(request.Key, 10, 64)
		if err != nil {
			break
		}
		if game, ok := allGames[intKey]; ok {
			response.Value = convertGameToMatchDetail(game)
			response.Ok = true
		}
	case "team":
	case "rankedData":
		intKey, err := strconv.ParseInt(request.Key, 10, 64)
		if err != nil {
			break
		}
		if data, ok := allRankedData[intKey]; ok {
			if data.ExpireTime > time.Now().Unix() {
				response.Value = data
				response.Ok = true
			} else {
				fmt.Printf("Cached ranked data is too old.")
			}
		}
	}

	request.Response <- response
}

func setupCache() {
	loadChampions(championCache)
	loadSummoners(summonerCache)
	allGames = make(map[int64]lapi.Game, 1)
	gamesBySummoner = make(map[int64][]lapi.Game, 1)
	allRankedData = make(map[int64]SummonerRankedData, 1)
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

func getExpireTime() int64 {
	return (time.Now().Add(time.Minute * cacheTimeoutMinutes)).Unix()
}
