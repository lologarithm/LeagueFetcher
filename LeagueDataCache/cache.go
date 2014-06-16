// Caching layer for use with the LeagueApi package.
package LeagueDataCache

import (
	"fmt"
	"strings"
	"time"

	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

const (
	championCache       = "ccache.json"
	summonerCache       = "scache.json"
	itemCache           = "icache.json"
	cacheTimeoutMinutes = 30
)

type PersistanceProvider interface {
	PutSummoner(lapi.Summoner) error
	GetSummoner(*lapi.Summoner) error
	GetSummonerByName(*lapi.Summoner) error
	GetSummoners([]int64) ([]lapi.Summoner, error)
	GetMatchDetail(MatchKey, *MatchDetail) error
	PutMatchDetail(MatchKey, MatchDetail) error
	PutMatchDetails(int64, []MatchDetail) error
	GetMatchDetails(int64) ([]MatchDetail, error)
	GetMatchHistory(int64) (MatchHistory, error)
}

// TODO: Maybe merge the request/response together?
type Request struct {
	Response chan Response
	Type     string
	Key      interface{}
}

type Response struct {
	Key   interface{}
	Value interface{}
	Type  string
	Ok    bool
}

var CacheRunning bool

// Cache dictionaries.
var allSummonersByName map[string]lapi.Summoner
var allSummonersById map[int64]lapi.Summoner
var allChampions map[int64]lapi.Champion
var allLeagues map[int64]lapi.League
var allTeams map[int64]lapi.Team
var allGames map[MatchKey]cachedMatchDetail
var gamesBySummoner map[int64][]cachedMatchDetail
var allRankedData map[int64]SummonerRankedData
var allItems *lapi.ItemList

// RunCache is the primary method. Start this as a goroutine and then use other public methods to fetch from this.
func RunCache(exit chan bool, get chan Request, put chan Response) {
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
			cacheSummoner(summoner)
		}
	case "champion":
		// For now the local cache will handle this.
	case "games":
		if matches, ok := resp.Value.(lapi.RecentGames); ok {
			for _, match := range matches.Games {
				md, cErr := convertGameToMatchDetail(match)
				if cErr != nil {
					fmt.Printf("Failed to convert game to MD: %s\n", cErr.Error())
					continue
				}
				cachedMatch := md.toCachedMatch(matches.SummonerId)
				key := MatchKey{MatchId: match.GameId, SummonerId: matches.SummonerId}
				putMatchInCache(key, cachedMatch)
			}
		}
	case "matchdetails":
		summonerId, ok := resp.Key.(int64)
		if !ok {
			return
		}
		if matches, ok := resp.Value.([]MatchDetail); ok {
			for _, match := range matches {
				cachedMatch := match.toCachedMatch(summonerId)
				key := MatchKey{MatchId: match.GameId, SummonerId: summonerId}
				putMatchInCache(key, cachedMatch)
			}
		}
	case "game":
		if key, ok := resp.Key.(MatchKey); ok {
			if game, ok := resp.Value.(MatchDetail); ok {
				putMatchInCache(key, game.toCachedMatch(key.SummonerId))
			}
		}
	case "rankedData":
		if data, ok := resp.Value.(SummonerRankedData); ok {
			data.ExpireTime = getExpireTime(false)
			allRankedData[data.Id] = data
		}
	}
}

func putMatchInCache(key MatchKey, cachedMatch cachedMatchDetail) {
	// Don't put the match into the list of games for the summoner if it's already cached.
	if _, ok := allGames[key]; !ok {
		if _, ok := gamesBySummoner[key.SummonerId]; !ok {
			gamesBySummoner[key.SummonerId] = []cachedMatchDetail{cachedMatch}
		} else {
			for ind, m := range gamesBySummoner[key.SummonerId] {
				if m.PlayedDate < cachedMatch.PlayedDate {
					gamesBySummoner[key.SummonerId] = append(append(gamesBySummoner[key.SummonerId][:ind], cachedMatch), gamesBySummoner[key.SummonerId][ind+1:]...)
					break
				} else if ind == len(gamesBySummoner[key.SummonerId])-1 {
					gamesBySummoner[key.SummonerId] = append(gamesBySummoner[key.SummonerId], cachedMatch)
					break
				}
			}
		}
	}
	// Always re-cache here for updated match time.
	allGames[key] = cachedMatch
}

func fetchCache(request Request) {
	response := &Response{Ok: false}
	switch request.Type {
	case "summoner":
		if key, ok := request.Key.(string); ok {
			key = NormalizeString(key)
			if summoner, ok := allSummonersByName[key]; ok {
				response.Ok = true
				response.Value = summoner
			}
		} else if iKey, ok := request.Key.(int64); ok {
			if summoner, ok := allSummonersById[iKey]; ok {
				response.Ok = true
				response.Value = summoner
			}
		}
	case "champion":
		if intKey, ok := request.Key.(int64); ok {
			if champ, ok := allChampions[intKey]; ok {
				response.Value = champ
				response.Ok = true
			}
		}
	case "item":
		if intKey, ok := request.Key.(int64); ok {
			if item, ok := allItems.ItemsById[intKey]; ok {
				response.Value = item
				response.Ok = true
			}
		}
	case "games":
		if intKey, ok := request.Key.(int64); ok {
			// If last fetch of game history is old, refetch game history.
			if games, ok := gamesBySummoner[intKey]; ok {
				checkExpireMatch := allGames[MatchKey{SummonerId: intKey, MatchId: games[0].GameId}]
				if len(games) > 0 && checkExpireMatch.CacheExpireDate > time.Now().Unix() {
					sliceEnd := 10
					if len(games) < 10 {
						sliceEnd = len(games)
					}
					matchHistory, fetchErr := convertCachedMatchDetailsToHistory(games[0:sliceEnd])
					if fetchErr == nil {
						response.Value = matchHistory
						response.Ok = true
					}
				}
			}
		}
	case "game":
		if key, ok := request.Key.(MatchKey); ok {
			if game, ok := allGames[key]; ok {
				gd, err := game.ToMatchDetail()
				if err == nil {
					response.Value = gd
					response.Ok = true
				} else {
					fmt.Printf("Convert to MatchDetail failed: %s", err.Error())
				}
			} else {
				fmt.Printf("GameKey not in allGames: %s\n", key.String())
			}
		}
	case "team":
	case "rankedData":
		if intKey, ok := request.Key.(int64); ok {
			if data, ok := allRankedData[intKey]; ok {
				if data.ExpireTime > time.Now().Unix() {
					response.Value = data
					response.Ok = true
				}
			}
		}
	}
	request.Response <- *response
}

func cacheSummoner(summoner lapi.Summoner) {
	allSummonersById[summoner.Id] = summoner
	key := NormalizeString(summoner.Name)
	allSummonersByName[key] = summoner
}

func SetupCache() {
	loadChampions(championCache)
	loadItems(itemCache)
	loadSummoners(summonerCache)

	if allGames == nil {
		allGames = make(map[MatchKey]cachedMatchDetail)
	}
	if gamesBySummoner == nil {
		gamesBySummoner = make(map[int64][]cachedMatchDetail)
	}
	if allRankedData == nil {
		allRankedData = make(map[int64]SummonerRankedData)
	}
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

func getExpireTime(inNano bool) int64 {
	if inNano {
		return (time.Now().Add(time.Minute * cacheTimeoutMinutes)).UnixNano()
	}
	return (time.Now().Add(time.Minute * cacheTimeoutMinutes)).Unix()
}

func NormalizeString(s string) string {
	s = strings.ToLower(s)
	return strings.Replace(s, " ", "", -1)
}
