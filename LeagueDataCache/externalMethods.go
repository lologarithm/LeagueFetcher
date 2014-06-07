package LeagueDataCache

import (
	"appengine"
	"appengine/urlfetch"
	"errors"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Public Functions that use channel to communicate with cache goroutine

// Fetch Champion from cache goroutine.
func GetChampion(id int64, get chan Request, c appengine.Context) (lapi.Champion, error) {
	if id <= 0 {
		// Return empty champion if there is no champion.
		return lapi.Champion{}, nil
	}
	result := make(chan Response, 1)
	cReq := Request{Type: "champion", Key: id, Response: result, Context: c}
	get <- cReq
	champResponse := <-result
	if champResponse.Ok {
		champ, _ := champResponse.Value.(lapi.Champion)
		return champ, nil
	}

	return lapi.Champion{}, CacheError{message: "Failed to retrieve champion."}
}

// Fetch Summoner from cache goroutine
func GetSummoner(name string, get chan Request, put chan Response, c appengine.Context, persist PersistanceProvider) (lapi.Summoner, error) {
	name = NormalizeString(name)
	value, getErr := goGet(Request{Type: "summoner", Key: name, Context: c, Persist: persist}, get)
	if getErr != nil {
		client := getClient(c)
		api := &lapi.LolFetcher{Get: client.Get, Log: c}
		// Summoner not cached. Retrieving from LAPI
		summoners, apiErr := api.GetSummonerByName(name)
		if apiErr != nil {
			// Need to handle different errors probably.
		}
		if s, gotOk := summoners[name]; gotOk {
			// Cache value before returning.
			goPut(s, "summoner", put, c, persist)
			return s, nil
		} else {
			return lapi.Summoner{}, CacheError{message: "Failed to retrieve summoner."}
		}
	}
	// Return the cached value!
	s, _ := value.(lapi.Summoner)
	return s, nil
}

// Fetch Summoner recent match history typed version
func GetSummonerMatchesSimple(id int64, get chan Request, put chan Response, c appengine.Context, persist PersistanceProvider) (MatchHistory, error) {
	value, getErr := goGet(Request{Type: "games", Key: id, Context: c}, get)
	if getErr != nil {
		client := getClient(c)
		api := &lapi.LolFetcher{Get: client.Get, Log: c}
		games, fetchErr := api.GetRecentMatches(id)
		if fetchErr != nil {
			return MatchHistory{}, fetchErr
		}
		if len(games.Games) > 0 {
			goPut(games, "games", put, c, persist)
			matches, fErr := convertGamesToMatchHistory(id, games.Games, func(id int64, api *lapi.LolFetcher) (lapi.Champion, error) {
				champ, fErr := GetChampion(id, get, c)
				return champ, fErr
			}, api)
			return matches, fErr
		}
		return MatchHistory{SummonerId: id}, CacheError{message: "No recent games found for summoner."}
	}
	lmh, _ := value.(MatchHistory)
	return lmh, nil
}

// Typed cache get for a match
func GetMatch(matchId int64, summonerId int64, get chan Request, put chan Response, c appengine.Context, persist PersistanceProvider) (MatchDetail, error) {
	value, getErr := goGet(Request{Type: "game", Key: MatchKey{MatchId: matchId, SummonerId: summonerId}, Context: c, Persist: persist}, get)

	client := getClient(c)
	api := &lapi.LolFetcher{Get: client.Get, Log: c}
	if getErr != nil {
		return MatchDetail{}, CacheError{message: "Failed to retrieve game."}
	}
	game, _ := value.(MatchDetail)
	missingIds := []int64{}
	for _, p := range game.FellowPlayers {
		if p.SummonerName == "" {
			missingIds = append(missingIds, p.SummonerId)
		}
	}

	if len(missingIds) > 0 {
		fetchedSummoners, apiErr := api.GetSummonersById(missingIds)
		if apiErr != nil {
			c.Infof("Failed to get summoners: %s", apiErr.Error())
			return game, apiErr
		}

		for _, value := range fetchedSummoners {
			// Cache this guy
			go goPut(value, "summoner", put, c, persist)
			// Put this summoner name where it belongs.
			// This is inefficient but there is never more than 10 players.. so meh.
			for i := 0; i < len(game.FellowPlayers); i++ {
				p := game.FellowPlayers[i]
				if p.SummonerId == value.Id {
					p.SummonerName = value.Name
					game.FellowPlayers[i] = p
					break
				}
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
	return game, nil
}

func GetSummonerRankedData(s lapi.Summoner, get chan Request, put chan Response, c appengine.Context) (srd SummonerRankedData, e error) {
	// 1. Check for cached data
	value, getErr := goGet(Request{Type: "rankedData", Key: s.Id, Context: c}, get)
	client := getClient(c)
	api := &lapi.LolFetcher{Get: client.Get, Log: c}
	if getErr != nil {
		srd.Summoner = s
		srd.RankedTeamLeagues = make(map[string]lapi.League)

		// Async get stats & leagues together.
		getStats := make(chan lapi.ApiAsyncResponse)
		getLeagues := make(chan lapi.ApiAsyncResponse)
		go api.GetSummonerRankedStatsAsync(s.Id, getStats)
		go api.GetSummonerLeaguesAsync(s.Id, getLeagues)
		responses := 0
		for {
			select {
			case statsAsync := <-getStats:
				if statsAsync.Error != nil {
					return srd, statsAsync.Error
				}
				if stats, ok := statsAsync.Value.(lapi.RankedStats); ok {
					for index, stat := range stats.Champions {
						cVal, _ := goGet(Request{Type: "champion", Key: stat.Id, Context: c}, get)
						champ, _ := cVal.(lapi.Champion)
						stat.ChampionName = champ.Name
						stats.Champions[index] = stat
					}
					srd.RankedStats = stats
				}
				// If we have both pieeces, return!
				if responses == 1 {
					goPut(srd, "rankedData", put, c, nil)
					return srd, nil
				}
				responses += 1
			case leaguesAsync := <-getLeagues:
				if leaguesAsync.Error != nil && leaguesAsync.Error.Error() == "Rate Limit Exceeded." {
					return srd, leaguesAsync.Error
				}
				if leagues, ok := leaguesAsync.Value.(map[string][]lapi.League); ok {
					for _, league := range leagues[strconv.FormatInt(s.Id, 10)] {
						if len(league.Entries) > 1 {
							// Possible bug state where they don't actually give us the stats.
							continue
						}
						if league.Queue == "RANKED_SOLO_5x5" {
							srd.Solo5sLeague = league
						} else if league.Queue == "RANKED_SOLO_3x3" {
							srd.Solo3sLeague = league
						} else if strings.Contains(league.Queue, "RANKED") {
							srd.RankedTeamLeagues[league.Entries[0].PlayerOrTeamId] = league
						}
					}
				}
				// If we have both pieeces, return!
				if responses == 1 {
					goPut(srd, "rankedData", put, c, nil)
					return srd, nil
				}
				responses += 1
			}
		}
		return srd, nil
	}

	if rankData, ok := value.(SummonerRankedData); ok {
		return rankData, nil
	}

	// If all else fails, return an empty object.
	return srd, errors.New("Failed to get stats.")
}

// Generic function to request data from cache.
func goGet(request Request, get chan Request) (interface{}, error) {
	cResponse := make(chan Response, 1)
	request.Response = cResponse
	get <- request
	respValue := <-cResponse
	if respValue.Ok {
		return respValue.Value, nil
	}
	return nil, CacheError{message: "Failed to fetch from cache!"}
}

// Generic function to cache a value
func goPut(value interface{}, valType string, put chan Response, c appengine.Context, pp PersistanceProvider) {
	r := Response{Value: value, Type: valType, Context: c, Persist: pp}
	put <- r
}

func getClient(ctx appengine.Context) *http.Client {
	if appengine.IsDevAppServer() {
		transport := http.Transport{}
		return &http.Client{Transport: &transport}
	}
	return urlfetch.Client(ctx)
}
