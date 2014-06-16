package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	lolCache "github.com/lologarithm/LeagueFetcher/LeagueDataCache"

	"appengine"
	"appengine/datastore"
)

type ServerConfig struct {
	ApiKey string
}

func loadConfig() {
	var cfg ServerConfig

	configData, readErr := ioutil.ReadFile("config.json")
	if readErr != nil {
		fmt.Printf("Error loading config: %s\n", readErr.Error())
		return
	}
	marshErr := json.Unmarshal(configData, &cfg)
	if marshErr != nil {
		fmt.Printf("Error parsing config: %s\n", marshErr.Error())
	}
	lapi.ApiKey = cfg.ApiKey
}

type endpointFunc func(http.ResponseWriter, *http.Request, chan lolCache.Request, chan lolCache.Response)

func init() {
	loadConfig()
	cacheGet := make(chan lolCache.Request, 10)
	cachePut := make(chan lolCache.Response, 10)
	exit := make(chan bool, 1)
	lolCache.CacheRunning = true
	lolCache.SetupCache()
	go lolCache.RunCache(exit, cacheGet, cachePut)

	http.HandleFunc("/", defaultHandler)
	// Wrap handlers with closure that passes in the channel for cache requests.

	// Page requests
	//	http.HandleFunc("/rankedStats", func(w http.ResponseWriter, req *http.Request) { handleRankedStats(w, req, cacheGet, cachePut) })

	// JSON data api

	// Static Data
	http.HandleFunc("/api/champion", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleChampion, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/api/item", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleItem, w, req, cacheGet, cachePut)
	})

	// Dynamic Data
	http.HandleFunc("/api/summoner/matchHistory", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleRecentMatches, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/api/summoner/rankedData", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleRankedData, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/api/summoner/LFScores", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleGetLFScore, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/api/match", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleMatchDetails, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/task/cacheGames", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleGameCache, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/task/calcStats", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleCalcStats, w, req, cacheGet, cachePut)
	})

	http.HandleFunc("/clean", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(cleanDatastore, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/log", handleLog)
}

// Wrapper function to time the endpoint call.
func timeEndpoint(endFunc endpointFunc, w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	st := time.Now().UnixNano()
	endFunc(w, r, cacheGet, cachePut)
	c.Infof("API Query (%s) Took: %.4fms\n", r.URL, (float64(time.Now().UnixNano()-st))/float64(1000000.0))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hit default handler: %s\n", r.URL)
	w.Write([]byte("Hello"))
}

func handleRecentMatches(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	summoner, fetchErr := lolCache.GetSummoner(html.UnescapeString(r.FormValue("name")), cacheGet, cachePut, c, &lolCache.MemcachePersistance{Context: c})

	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}
	matches, fetchErr := lolCache.GetSummonerMatchesSimple(summoner.Id, cacheGet, cachePut, c, &lolCache.MemcachePersistance{Context: c})

	// Now send tasks to start caching all games from other summoner perspecives.ServerConfig
	//for _, match := range matches.Games {
	//	t := taskqueue.NewPOSTTask("/task/cacheGames", map[string][]string{"matchId": {strconv.FormatInt(match.GameId, 10)}, "summonerId": {strconv.FormatInt(summoner.Id, 10)}})
	//	taskqueue.Add(c, t, "")
	//}

	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}
	writeJson(w, matches, c)
}

func handleMatchDetails(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	matchId, intErr := strconv.ParseInt(r.FormValue("matchId"), 10, 64)
	if intErr != nil {
		returnErrJson(intErr, w, c)
		return
	}
	summonerId, intErr := strconv.ParseInt(r.FormValue("summonerId"), 10, 64)
	if intErr != nil {
		returnErrJson(intErr, w, c)
		return
	}

	match, fetchErr := lolCache.GetMatch(matchId, summonerId, cacheGet, cachePut, c, &lolCache.MemcachePersistance{Context: c})
	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}

	writeJson(w, match, c)
}

func handleChampion(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	champId, parseErr := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if parseErr != nil {
		returnEmptyJson(w)
		return
	}
	champ, fetchErr := lolCache.GetChampion(champId, cacheGet, c)
	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}
	writeJson(w, champ, c)
}

func handleItem(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	itemId, parseErr := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if parseErr != nil {
		returnEmptyJson(w)
		return
	}

	item, fetchErr := lolCache.GetItem(itemId, cacheGet, cachePut, c)
	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}
	writeJson(w, item, c)
}

func handleRankedData(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	summoner, fetchErr := lolCache.GetSummoner(html.UnescapeString(r.FormValue("name")), cacheGet, cachePut, c, &lolCache.MemcachePersistance{Context: c})
	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}

	data, fetchErr := lolCache.GetSummonerRankedData(summoner, cacheGet, cachePut, c)
	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}
	writeJson(w, data, c)
}

func handleGameCache(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	matchId, intErr := strconv.ParseInt(r.FormValue("matchId"), 10, 64)
	if intErr != nil {
		returnErrJson(intErr, w, c)
		return
	}
	summonerId, intErr := strconv.ParseInt(r.FormValue("summonerId"), 10, 64)
	if intErr != nil {
		returnErrJson(intErr, w, c)
		return
	}

	lolCache.CacheMatch(matchId, summonerId, cacheGet, cachePut, c, &lolCache.MemcachePersistance{Context: c})
}

func handleGetLFScore(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	c := appengine.NewContext(r)
	summonerId, intErr := strconv.ParseInt(r.FormValue("summonerId"), 10, 64)
	if intErr != nil {
		returnErrJson(intErr, w, c)
		return
	}

	scores, fetchErr := lolCache.GetLFScores(summonerId, cacheGet, cachePut, c, &lolCache.MemcachePersistance{Context: c})
	if fetchErr != nil {
		returnErrJson(fetchErr, w, c)
		return
	}
	writeJson(w, scores, c)
}

func handleCalcStats(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
}

func cleanDatastore(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	if r.FormValue("p") != "ben" {
		returnEmptyJson(w)
		return
	}

	c := appengine.NewContext(r)
	var keys []*datastore.Key
	keys, err := datastore.NewQuery("Match").Filter("IntIndex = ", nil).KeysOnly().GetAll(c, keys)
	if err != nil {
		returnErrJson(err, w, c)
	}
	datastore.DeleteMulti(c, keys)
}

func returnEmptyJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func returnErrJson(e error, w http.ResponseWriter, context appengine.Context) {
	w.Header().Set("Content-Type", "application/json")
	msg := fmt.Sprintf("{\"error\": \"%s\"}", e.Error())
	context.Infof("Returning Error: %s", msg)
	w.Write([]byte(msg))
}

func writeJson(w http.ResponseWriter, data interface{}, context appengine.Context) {
	w.Header().Set("Content-Type", "application/json")
	dataJson, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		returnErrJson(jsonErr, w, context)
		return
	}
	w.Write(dataJson)
}

func handleLog(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	ctx.Infof(req.FormValue("msg"))
}

//func handleRankedStats(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
//	summoner, err := lolCache.GetSummoner(r.FormValue("name"), cacheGet, cachePut)
//	w.Write([]byte("<html><body><pre>"))
//	if err == nil {
//		w.Write(GetRankedStats(summoner, cacheGet, cachePut))
//	}
//	w.Write([]byte("</pre></body></html>"))
//}
//func GetSummonerSummary(s lapi.Summoner) []byte {
//	var buffer bytes.Buffer

//	summary := lapi.GetSummonerSummaryStats(s.Id)
//	for _, stats := range summary.PlayerStatSummaries {
//		if strings.Contains(stats.PlayerStatSummaryType, "Ranked") {
//			buffer.WriteString(fmt.Sprintf("Game Type: %s\n  Wins: %d\n  Losses: %d\n  Kills: %d\n  Assists: %d\n", stats.PlayerStatSummaryType, stats.Wins, stats.Losses, stats.AggregatedStats.TotalChampionKills, stats.AggregatedStats.TotalAssists))
//		} else {
//			buffer.WriteString(fmt.Sprintf("Game Type: %s\n  Wins: %d\n  Kills: %d\n  Assists: %d\n", stats.PlayerStatSummaryType, stats.Wins, stats.AggregatedStats.TotalChampionKills, stats.AggregatedStats.TotalAssists))
//		}
//	}

//	return buffer.Bytes()
//}

//func GetRankedStats(s lapi.Summoner, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) []byte {
//	var buffer bytes.Buffer

//	stats := lapi.GetSummonerRankedStats(s.Id)
//	leagues := lapi.GetSummonerLeagues(s.Id)
//	soloTierDiv := "Cardboard V"
//	teamTier := "Cardboard"
//	teamDivision := "V"
//	for _, league := range leagues[strconv.FormatInt(s.Id, 10)] {
//		if league.Queue == "RANKED_SOLO_5x5" {
//			soloTierDiv = fmt.Sprintf("%s %s", league.Tier, league.Entries[0].Division)
//		} else if league.Queue == "RANKED_TEAM_5x5" {
//			if lapi.CompareRanked(league.Tier, league.Entries[0].Division, teamTier, teamDivision) == 1 {
//				teamTier = league.Tier
//				teamDivision = league.Entries[0].Division
//			}
//		}
//	}
//	buffer.WriteString(fmt.Sprintf("%s:\n  Solo Queue League: %s\n  Best Ranked 5's League: %s %s\nChampion Stats:\n", s.Name, soloTierDiv, teamTier, teamDivision))
//	for _, champStats := range stats.Champions {
//		if champStats.Id > 0 {
//			champ, champErr := lolCache.GetChampion(champStats.Id, cacheGet)
//			if champErr != nil {
//				// Not much I can do here...
//			}
//			buffer.WriteString(fmt.Sprintf("  Champ: %s,", champ.Name))
//		} else {
//			buffer.WriteString(fmt.Sprintf("\n  Champion Totals: "))
//		}
//		buffer.WriteString(fmt.Sprintf(
//			" Wins: %d, Losses: %d, Kills: %d, Assists: %d, Deaths: %d\n",
//			champStats.Stats.TotalSessionsWon, champStats.Stats.TotalSessionsLost, champStats.Stats.TotalChampionKills, champStats.Stats.TotalAssists, champStats.Stats.TotalDeathsPerSession))
//	}

//	return buffer.Bytes()
//}
