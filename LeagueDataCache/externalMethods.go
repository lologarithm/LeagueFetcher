package LeagueDataCache

import (
	"appengine"
	"appengine/urlfetch"
	"errors"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"math/rand"
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
	cReq := Request{Type: "champion", Key: id, Response: result}
	get <- cReq
	champResponse := <-result
	if champResponse.Ok {
		champ, _ := champResponse.Value.(lapi.Champion)
		return champ, nil
	}

	return lapi.Champion{}, CacheError{message: "Failed to retrieve champion."}
}

// Fetch Item from cache goroutine.
func GetItem(id int64, get chan Request, put chan Response, c appengine.Context) (lapi.Item, error) {
	if id <= 0 {
		// Return empty champion if there is no champion.
		return lapi.Item{}, nil
	}
	itemResponse, fErr := goGet(Request{Type: "item", Key: id}, get)
	if fErr != nil {
		// TODO: Fetch item from LAPI and cache
		return lapi.Item{}, CacheError{message: "Failed to retrieve item."}
	}
	item, _ := itemResponse.(lapi.Item)
	return item, nil

}

// Fetch Summoner from cache goroutine
func GetSummoner(name string, get chan Request, put chan Response, c appengine.Context, persist PersistanceProvider) (lapi.Summoner, error) {
	name = NormalizeString(name)

	// 1. Check Local Cache
	value, getErr := goGet(Request{Type: "summoner", Key: name}, get)
	if getErr != nil {
		// 2. Try to fetch from DB
		sRef := &lapi.Summoner{NormalizedName: name}
		dbErr := persist.GetSummonerByName(sRef)
		if dbErr != nil {
			client := getClient(c)
			api := &lapi.LolFetcher{Get: client.Get, Log: c}
			// Summoner not cached. Retrieving from LAPI
			summoners, apiErr := api.GetSummonerByName(name)
			if apiErr != nil {
				// Need to handle different errors probably.
				return lapi.Summoner{}, apiErr
			}
			if s, gotOk := summoners[name]; gotOk {
				// Cache value before returning.
				goPut(s, "summoner", put)
				fail := persist.PutSummoner(s)
				if fail != nil {
					c.Warningf("Failed to store summoner: %s", fail.Error())
				}
				return s, nil
			} else {
				return lapi.Summoner{}, CacheError{message: "Failed to retrieve summoner."}
			}
		} else {
			// Locally Cache
			goPut(*sRef, "summoner", put)
			// Return value
			return *sRef, nil
		}
	}

	// Return the cached value!
	s, _ := value.(lapi.Summoner)
	return s, nil
}

// Fetch Summoner recent match history typed version
func GetSummonerMatchesSimple(id int64, get chan Request, put chan Response, c appengine.Context, persist PersistanceProvider) (MatchHistory, error) {
	// 1. Try to fetch from local cache
	value, getErr := goGet(Request{Type: "games", Key: id}, get)
	if getErr != nil {
		c.Infof("Failed to get from local cache. Checking persistance.")
		games := lapi.RecentGames{SummonerId: id, Games: []lapi.Game{}}
		// Now try to get from persistance (db)
		gameList, persistErr := persist.GetMatchesByIndex(id)
		shouldPersist := false
		if persistErr != nil {
			shouldPersist = true
			if persistErr != nil {
				c.Warningf("FETCH FROM PERSIST FAILED: %s", persistErr)
			}
			client := getClient(c)
			api := &lapi.LolFetcher{Get: client.Get, Log: c}
			// I guess now we try to fetch from lapi
			gotGames, fetchErr := api.GetRecentMatches(id)
			if fetchErr != nil {
				// Don't return a 404 as an error. It means we found nothing
				if fetchErr.Code != 404 {
					return MatchHistory{}, nil
				}
				// Return err if we were unable to get from lapi
				return MatchHistory{}, fetchErr
			}
			games.Games = gotGames.Games
		} else {
			games.Games = gameList
		}
		// Recache locally
		c.Infof("Matches not stored in local cache. Storing")
		goPut(games, "games", put)
		// Now we need convert.
		if len(games.Games) > 0 {
			// Now convert to our local simple game format.
			matches, fErr := convertGamesToMatchHistory(id, games.Games)
			// Cache in 'persistance' if there was a persist error
			if shouldPersist {
				c.Infof("Matches not stored in persistance. Storing")
				persistGames(games, id, persist)
			}
			// Return the matches!
			return matches, fErr
		}
		return MatchHistory{SummonerId: id}, CacheError{message: "No recent games found for summoner."}
	}
	lmh, _ := value.(MatchHistory)
	return lmh, nil
}

// Typed cache get for a match
func GetMatch(matchId int64, summonerId int64, get chan Request, put chan Response, c appengine.Context, persist PersistanceProvider) (MatchDetail, error) {
	// 1. Try to fetch from local cache.
	mKey := MatchKey{MatchId: matchId, SummonerId: summonerId}
	value, getErr := goGet(Request{Type: "game", Key: mKey}, get)

	client := getClient(c)
	api := &lapi.LolFetcher{Get: client.Get, Log: c}
	if getErr != nil {
		// Now try to get from persistance
		var persistGame *lapi.Game
		persistErr := persist.GetObject("Match", mKey.String(), &persistGame)
		if persistErr != nil {
			return MatchDetail{}, persistErr
		} else {
			gameDetail, fErr := convertGameToMatchDetail(*persistGame)
			if fErr == nil {
				value = gameDetail
			}
		}
	}

	game, _ := value.(MatchDetail)
	missingIds := []int64{}
	for _, p := range game.FellowPlayers {
		if p.SummonerName == "" {
			missingIds = append(missingIds, p.SummonerId)
		}
	}

	if len(missingIds) > 0 {
		// First check local db.
		pSummoners, dbErr := persist.GetSummoners(missingIds)
		for _, value := range pSummoners {
			goPut(value, "summoner", put)
			// Put the name into the list
			for i := 0; i < len(game.FellowPlayers); i++ {
				p := game.FellowPlayers[i]
				if p.SummonerId == value.Id {
					p.SummonerName = value.Name
					game.FellowPlayers[i] = p
					break
				}
			}
			if dbErr != nil {
				// Now clean up the missingIds so we know which ones we are missing.
				for index, id := range missingIds {
					if id == value.Id {
						// Remove the ID from the list of missing.
						missingIds = append(missingIds[:index], missingIds[index+1:]...)
						break
					}
				}
			}
		}
		// Only fetch from API if not all summoners were returned from persistance.
		if dbErr != nil {
			fetchedSummoners, apiErr := api.GetSummonersById(missingIds)
			if apiErr != nil {
				return game, apiErr
			}
			for _, value := range fetchedSummoners {
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
		}
	}
	return game, nil
}

// Fetch and Cache all summoners and their perspective on the game.
func CacheMatch(matchId int64, summonerId int64, get chan Request, put chan Response, c appengine.Context, persist PersistanceProvider) {
	// Sleep random amount of time so we don't get too much overlapping
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	mKey := MatchKey{MatchId: matchId, SummonerId: summonerId}
	c.Infof("Fetching details of: %v to cache other summoners.", mKey)
	value, getErr := goGet(Request{Type: "game", Key: mKey}, get)
	var game MatchDetail
	if getErr != nil {
		var persistGame *lapi.Game
		persistErr := persist.GetObject("Match", mKey.String(), &persistGame)
		if persistErr != nil {
			// Couldn't find game to cache
			return
		}
		gameDetail, fErr := convertGameToMatchDetail(*persistGame)
		if fErr == nil {
			game = gameDetail
		}
	} else {
		game, _ = value.(MatchDetail)
	}

	missingIds := []int64{}
	for _, p := range game.FellowPlayers {
		if p.SummonerName == "" {
			missingIds = append(missingIds, p.SummonerId)
		}
	}

	client := getClient(c)
	api := &lapi.LolFetcher{Get: client.Get, Log: c}
	if len(missingIds) > 0 {
		// First check local db.
		pSummoners, dbErr := persist.GetSummoners(missingIds)
		for _, value := range pSummoners {
			goPut(value, "summoner", put)
			for index, id := range missingIds {
				if id == value.Id {
					// Remove the ID from the list of missing.
					missingIds = append(missingIds[:index], missingIds[index+1:]...)
					break
				}
			}
		}
		if dbErr != nil {
			aSummoners, apiErr := api.GetSummonersById(missingIds)
			if apiErr != nil {
				// Can't seem to make it work. Bail
				return
			}
			for _, value := range aSummoners {
				persist.PutSummoner(value)
				goPut(value, "summoner", put)
			}
		}
	}

	//c.Infof("Now fetching match history for each player in game.")
	//for _, fplay := range game.FellowPlayers {
	//	// Sleep to slow down requests.
	//	time.Sleep(time.Second)
	//	_, getErr := goGet(Request{Type: "game", Key: MatchKey{MatchId: matchId, SummonerId: fplay.SummonerId}, Context: c}, get)
	//	if getErr != nil {
	//		games, fetchErr := api.GetRecentMatches(fplay.SummonerId)
	//		if fetchErr != nil {
	//			continue
	//		}
	//		if len(games.Games) > 0 {
	//			goPut(games, "games", put, c)
	//			persistGames(games, fplay.SummonerId, persist)
	//		}
	//	}
	//}
}

func persistGames(games lapi.RecentGames, id int64, persist PersistanceProvider) {
	gameKeys := make([]string, len(games.Games))
	indexes := make([]int64, len(games.Games))
	for ind, g := range games.Games {
		gameKeys[ind] = (MatchKey{MatchId: g.GameId, SummonerId: id}).String()
		indexes[ind] = id
	}
	things := make([]interface{}, len(games.Games))
	for i, v := range games.Games {
		things[i] = interface{}(v)
	}
	persist.PutObjects("Match", gameKeys, indexes, things)
}

func GetSummonerRankedData(s lapi.Summoner, get chan Request, put chan Response, c appengine.Context) (srd SummonerRankedData, e error) {
	// 1. Check for cached data
	value, getErr := goGet(Request{Type: "rankedData", Key: s.Id}, get)
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
					if statsAsync.Error.Code != 404 {
						c.Infof("Error: %s", statsAsync.Error.Error())
						return srd, statsAsync.Error
					}
				}
				if stats, ok := statsAsync.Value.(lapi.RankedStats); ok {
					for index, stat := range stats.Champions {
						cVal, _ := goGet(Request{Type: "champion", Key: stat.Id}, get)
						champ, _ := cVal.(lapi.Champion)
						stat.ChampionName = champ.Name
						if stat.Id > 0 {
							stat.ChampionImage = champ.Image.GetImageURL()
						}
						stats.Champions[index] = stat
					}
					srd.RankedStats = stats
				}
				// If we have both pieeces, return!
				if responses == 1 {
					goPut(srd, "rankedData", put)
					return srd, nil
				}
				responses += 1
			case leaguesAsync := <-getLeagues:
				if leaguesAsync.Error != nil {
					// Don't return a 404 as an error. It means we found nothing
					if leaguesAsync.Error.Code != 404 {
						c.Infof("Error: %s", leaguesAsync.Error.Error())
						return srd, leaguesAsync.Error
					}
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
					goPut(srd, "rankedData", put)
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
func goPut(value interface{}, valType string, put chan Response) {
	r := Response{Value: value, Type: valType}
	put <- r
}

func getClient(ctx appengine.Context) *http.Client {
	if appengine.IsDevAppServer() {
		transport := http.Transport{}
		return &http.Client{Transport: &transport}
	}
	return urlfetch.Client(ctx)
}
