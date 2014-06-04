package main

import (
	"encoding/json"
	"fmt"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	lolCache "github.com/lologarithm/LeagueFetcher/LeagueDataCache"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

func main() {
	loadConfig()
	cacheGet := make(chan lolCache.Request, 10)
	cachePut := make(chan lolCache.Response, 10)
	exit := make(chan bool, 1)
	go lolCache.RunCache(exit, cacheGet, cachePut)

	http.HandleFunc("/", defaultHandler)
	// Wrap handlers with closure that passes in the channel for cache requests.

	// Page requests
	//	http.HandleFunc("/rankedStats", func(w http.ResponseWriter, req *http.Request) { handleRankedStats(w, req, cacheGet, cachePut) })

	// JSON data api
	http.HandleFunc("/api/champion", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleChampion, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/api/summoner/matchHistory", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleRecentMatches, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/api/summoner/rankedData", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleRankedData, w, req, cacheGet, cachePut)
	})
	http.HandleFunc("/api/match", func(w http.ResponseWriter, req *http.Request) {
		timeEndpoint(handleMatchDetails, w, req, cacheGet, cachePut)
	})

	http.ListenAndServe(":9000", nil)
}

// Wrapper function to time the endpoint call.
func timeEndpoint(endFunc endpointFunc, w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	st := time.Now().UnixNano()
	endFunc(w, r, cacheGet, cachePut)
	fmt.Printf("Request (%s) Took: %.4fms\n", r.URL, (float64(time.Now().UnixNano()-st))/float64(1000000.0))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hit default handler: %s\n", r.URL)
	w.Write([]byte("Hello"))
}

func handleRecentMatches(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	summoner, err := lolCache.GetSummoner(r.FormValue("name"), cacheGet, cachePut)
	if err != nil {
		returnEmptyJson(w)
		return
	}
	matches, _ := lolCache.GetSummonerMatchesSimple(summoner.Id, cacheGet, cachePut)
	writeJson(w, matches)
}

func handleMatchDetails(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	intKey, intErr := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if intErr != nil {
		returnEmptyJson(w)
		return
	}
	match, _ := lolCache.GetMatch(intKey, cacheGet, cachePut)
	writeJson(w, match)
}

func handleChampion(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	fmt.Printf("R: %v", r)
	champId, parseErr := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if parseErr != nil {
		returnEmptyJson(w)
		return
	}
	champ, err := lolCache.GetChampion(champId, cacheGet)
	if err != nil {
		returnEmptyJson(w)
		return
	}
	writeJson(w, champ)
}

func handleRankedData(w http.ResponseWriter, r *http.Request, cacheGet chan lolCache.Request, cachePut chan lolCache.Response) {
	summoner, err := lolCache.GetSummoner(r.FormValue("name"), cacheGet, cachePut)
	if err != nil {
		returnEmptyJson(w)
		return
	}

	data := lolCache.GetSummonerRankedData(summoner, cacheGet, cachePut)
	writeJson(w, data)
}

func returnEmptyJson(w http.ResponseWriter) {
	w.Write([]byte("{}"))
}

func writeJson(w http.ResponseWriter, data interface{}) {
	dataJson, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		returnEmptyJson(w)
		return
	}
	w.Write(dataJson)
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
