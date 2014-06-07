// Caching layer for use with the LeagueApi package.
package LeagueDataCache

import (
	"appengine"
	"appengine/urlfetch"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"strings"
	"time"
)

const (
	championCache       = "ccache.json"
	summonerCache       = "scache.json"
	cacheTimeoutMinutes = 30
)

type PersistanceProvider interface {
	PutSummoner(lapi.Summoner) error
	GetSummoner(*lapi.Summoner) error
	GetSummonerByName(*lapi.Summoner) error
	//PutMatch(lapi.Game)
	PutObject(string, string, interface{}) error
	GetObject(string, string, interface{}) error
}

// TODO: Maybe merge the request/response together?
type Request struct {
	Response chan Response
	Type     string
	Key      interface{}
	Context  appengine.Context
	Persist  PersistanceProvider
}

type Response struct {
	Context appengine.Context
	Key     interface{}
	Value   interface{}
	Type    string
	Ok      bool
	Persist PersistanceProvider
}

var CacheRunning bool

// Cache dictionaries.
var allSummonersByName map[string]lapi.Summoner
var allSummonersById map[int64]lapi.Summoner
var allChampions map[int64]lapi.Champion
var allLeagues map[int64]lapi.League
var allTeams map[int64]lapi.Team
var allGames map[MatchKey]lapi.Game
var gamesBySummoner map[int64][]lapi.Game
var allRankedData map[int64]SummonerRankedData

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
			//resp.Persist.PutObject("Summoner", key, summoner)
			resp.Persist.PutSummoner(summoner)
		}
	case "champion":
		// For now the local cache will handle this.
	case "games":
		if matches, ok := resp.Value.(lapi.RecentGames); ok {
			for _, match := range matches.Games {
				match.ExpireTime = getExpireTime()
				key := MatchKey{MatchId: match.GameId, SummonerId: matches.SummonerId}
				// Don't put the match into the list of games for the summoner if it's already cached.
				if _, ok := allGames[key]; !ok {
					if _, ok := gamesBySummoner[matches.SummonerId]; !ok {
						gamesBySummoner[matches.SummonerId] = []lapi.Game{}
					}
					gamesBySummoner[matches.SummonerId] = append([]lapi.Game{match}, gamesBySummoner[matches.SummonerId]...)
				}
				// Always re-cache here for updated match time.
				allGames[key] = match
				go resp.Persist.PutObject("Match", key.String(), match)
			}
		} else {
			resp.Context.Infof("Cache: Failed to get game, not a game key!")
		}
	case "game":
		// Fix this up later if we ever use it.
		if key, ok := resp.Key.(MatchKey); ok {
			if game, ok := resp.Value.(lapi.Game); ok {
				allGames[key] = game
				resp.Persist.PutObject("Match", key.String(), game)
			}
		}
	case "rankedData":
		if data, ok := resp.Value.(SummonerRankedData); ok {
			data.ExpireTime = getExpireTime()
			allRankedData[data.Id] = data
		}
	}
}

func fetchCache(request Request) {
	client := urlfetch.Client(request.Context)
	api := &lapi.LolFetcher{Get: client.Get, Log: request.Context}
	wrappedFetch(request, api)
}

func wrappedFetch(request Request, api *lapi.LolFetcher) {
	response := &Response{Ok: false}
	switch request.Type {
	case "summoner":
		if key, ok := request.Key.(string); ok {
			key = NormalizeString(key)
			if summoner, ok := allSummonersByName[key]; ok {
				response.Ok = true
				response.Value = summoner
			} else {
				persistSummoner := &lapi.Summoner{Name: key}
				//persistErr := request.Persist.GetObject("Summoner", key, &persistSummoner)
				persistErr := request.Persist.GetSummonerByName(persistSummoner)
				if persistErr == nil {
					cacheSummoner(*persistSummoner)
					response.Value = *persistSummoner
					response.Ok = true
				}
			}
		}
	case "champion":
		if intKey, ok := request.Key.(int64); ok {
			champ, fErr := fetchAndCacheChampion(intKey, api)
			if fErr == nil {
				response.Value = champ
				response.Ok = true
			}
		}
	case "games":
		if intKey, ok := request.Key.(int64); ok {
			// If last fetch of game history is old, refetch game history.
			if games, ok := gamesBySummoner[intKey]; ok {
				checkExpireMatch := allGames[MatchKey{SummonerId: intKey, MatchId: games[0].GameId}]
				if len(games) > 0 && checkExpireMatch.ExpireTime > time.Now().Unix() {
					sliceEnd := 10
					if len(games) < 10 {
						sliceEnd = len(games)
					}
					matchHistory, fetchErr := convertGamesToMatchHistory(intKey, games[0:sliceEnd], fetchAndCacheChampion, api)
					if fetchErr == nil {
						response.Value = matchHistory
						response.Ok = true
					}
				} else {
					request.Context.Warningf("Cached games are old or no games found.")
				}
			} else {
				request.Context.Warningf("Cache: No games found for that id.")
			}
		} else {
			request.Context.Warningf("Cache: Not int key for fetching games.")
		}
	case "game":
		if key, ok := request.Key.(MatchKey); ok {
			if game, ok := allGames[key]; ok {
				gameDetail, fErr := convertGameToMatchDetail(game, api)
				if fErr == nil {
					response.Value = gameDetail
					response.Ok = true
				}
			} else {
				var persistGame *lapi.Game
				persistErr := request.Persist.GetObject("Match", key.String(), &persistGame)
				if persistErr != nil {
					request.Context.Warningf("Cache: Failed to get game: %v", key)
				} else {
					gameDetail, fErr := convertGameToMatchDetail(*persistGame, api)
					if fErr == nil {
						response.Value = gameDetail
						response.Ok = true
					}
				}
			}
		} else {
			request.Context.Warningf("Cache: Failed to get game, not a game key!")
		}
	case "team":
	case "rankedData":
		if intKey, ok := request.Key.(int64); ok {
			if data, ok := allRankedData[intKey]; ok {
				if data.ExpireTime > time.Now().Unix() {
					response.Value = data
					response.Ok = true
				} else {
					request.Context.Warningf("Cached ranked data is too old.")
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
	loadSummoners(summonerCache)
	if allGames == nil {
		allGames = make(map[MatchKey]lapi.Game)
	}
	if gamesBySummoner == nil {
		gamesBySummoner = make(map[int64][]lapi.Game)
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

func getExpireTime() int64 {
	return (time.Now().Add(time.Minute * cacheTimeoutMinutes)).Unix()
}

func NormalizeString(s string) string {
	s = strings.ToLower(s)
	return strings.Replace(s, " ", "", -1)
}
