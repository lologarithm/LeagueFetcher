package LeagueDataCache

import (
	"fmt"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"strconv"
	"strings"
)

// Public Functions that use channel to communicate with cache goroutine

// Fetch Champion from cache goroutine.
func GetChampion(id int64, get chan Request) (c lapi.Champion, e error) {
	result := make(chan Response, 1)
	cReq := Request{Type: "champion", Key: fmt.Sprintf("%d", id), Response: result}
	get <- cReq
	champResponse := <-result
	if champResponse.Ok {
		champ, _ := champResponse.Value.(lapi.Champion)
		return champ, nil
	}

	return lapi.Champion{}, CacheError{message: "Failed to retrieve champion."}
}

// Fetch Summoner from cache goroutine
func GetSummoner(name string, get chan Request, put chan Response) (lapi.Summoner, error) {
	value, getErr := goGet(Request{Type: "summoner", Key: name}, get)
	if getErr != nil {
		// Summoner not cached. Retrieving from LAPI
		summoners := lapi.GetSummonerByName(name)
		if s, gotOk := summoners[name]; gotOk {
			// Cache value before returning.
			goPut(s, "summoner", put)
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
func GetSummonerMatchesSimple(id int64, get chan Request, put chan Response) (MatchHistory, error) {
	value, getErr := goGet(Request{Type: "games", Key: fmt.Sprintf("%d", id)}, get)
	if getErr != nil {
		games := lapi.GetRecentMatches(id)
		if len(games.Games) > 0 {
			goPut(games, "games", put)
			matches := convertGamesToMatchHistory(id, games.Games, func(id int64) lapi.Champion {
				champ, _ := GetChampion(id, get)
				return champ
			})
			return matches, nil
		}
		return MatchHistory{SummonerId: id}, CacheError{message: "No recent games found for summoner."}
	}
	lmh, _ := value.(MatchHistory)
	return lmh, nil
}

// Typed cache get for a match
func GetMatch(id int64, get chan Request, put chan Response) (MatchDetail, error) {
	value, getErr := goGet(Request{Type: "game", Key: fmt.Sprintf("%d", id)}, get)
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
		fetchedSummoners := lapi.GetSummonersById(missingIds)
		for _, value := range fetchedSummoners {
			// Cache this guy
			goPut(value, "summoner", put)
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
	return game, nil
}

func GetSummonerRankedData(s lapi.Summoner, get chan Request, put chan Response) (srd SummonerRankedData) {
	// 1. Check for cached data
	value, getErr := goGet(Request{Type: "rankedData", Key: fmt.Sprintf("%d", s.Id)}, get)
	if getErr != nil {
		srd.Summoner = s
		// 1. Get RankedStats
		stats := lapi.GetSummonerRankedStats(s.Id)
		for index, stat := range stats.Champions {
			cVal, _ := goGet(Request{Type: "champion", Key: fmt.Sprintf("%d", stat.Id)}, get)
			champ, _ := cVal.(lapi.Champion)
			stat.ChampionName = champ.Name
			stats.Champions[index] = stat
		}
		srd.RankedStats = stats
		// 2. Get Teams... for now we don't do this.

		// 3. Get Leagues
		leagues := lapi.GetSummonerLeagues(s.Id)
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

		goPut(srd, "rankedData", put)
		return srd
	}

	if rankData, ok := value.(SummonerRankedData); ok {
		return rankData
	}

	// If all else fails, return an empty object.
	return srd
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
