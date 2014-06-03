package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
	apiKey = cfg.ApiKey
}

func main() {
	loadConfig()
	cacheRequests := make(chan CacheRequest, 10)
	exit := make(chan bool, 1)
	go RunCache(exit, cacheRequests)

	http.HandleFunc("/", defaultHandler)
	// Wrap handlers with closure that passes in the channel for cache requests.

	// Page requests
	http.HandleFunc("/rankedStats", func(w http.ResponseWriter, req *http.Request) { handleRankedStats(w, req, cacheRequests) })

	// JSON data api
	http.HandleFunc("/api/champion", func(w http.ResponseWriter, req *http.Request) { handleChampion(w, req, cacheRequests) })
	http.HandleFunc("/api/summoner/matchHistory", func(w http.ResponseWriter, req *http.Request) { handleRecentMatches(w, req, cacheRequests) })

	http.ListenAndServe(":9000", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

}

func handleRecentMatches(w http.ResponseWriter, r *http.Request, cacheRequests chan CacheRequest) {
	summoner, err := GetSummoner(r.FormValue("name"), cacheRequests)
	if err != nil {
		returnEmptyJson(w)
		return
	}
	matches := getRecentMatches(summoner.Id)
	matchJson, jsonErr := json.Marshal(matches)
	if jsonErr != nil {
		returnEmptyJson(w)
		return
	}
	w.Write(matchJson)
}

func handleChampion(w http.ResponseWriter, r *http.Request, cacheRequests chan CacheRequest) {
	fmt.Printf("R: %v", r)
	champId, parseErr := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if parseErr != nil {
		returnEmptyJson(w)
		return
	}
	champ, err := GetChampion(champId, cacheRequests)
	if err != nil {
		returnEmptyJson(w)
		return
	}
	champJson, jsonErr := json.Marshal(champ)
	if jsonErr != nil {
		returnEmptyJson(w)
		return
	}

	w.Write(champJson)
}

func returnEmptyJson(w http.ResponseWriter) {
	w.Write([]byte("{}"))
}

func handleRankedStats(w http.ResponseWriter, r *http.Request, cacheRequests chan CacheRequest) {
	summoner, err := GetSummoner(r.FormValue("name"), cacheRequests)
	w.Write([]byte("<html><body><pre>"))
	if err == nil {
		w.Write(GetRankedStats(summoner, cacheRequests))
	}
	w.Write([]byte("</pre></body></html>"))
}

func GetSummonerSummary(s Summoner) []byte {
	var buffer bytes.Buffer

	summary := getSummonerSummaryStats(s.Id)
	for _, stats := range summary.PlayerStatSummaries {
		if strings.Contains(stats.PlayerStatSummaryType, "Ranked") {
			buffer.WriteString(fmt.Sprintf("Game Type: %s\n  Wins: %d\n  Losses: %d\n  Kills: %d\n  Assists: %d\n", stats.PlayerStatSummaryType, stats.Wins, stats.Losses, stats.AggregatedStats.TotalChampionKills, stats.AggregatedStats.TotalAssists))
		} else {
			buffer.WriteString(fmt.Sprintf("Game Type: %s\n  Wins: %d\n  Kills: %d\n  Assists: %d\n", stats.PlayerStatSummaryType, stats.Wins, stats.AggregatedStats.TotalChampionKills, stats.AggregatedStats.TotalAssists))
		}
	}

	return buffer.Bytes()
}

func GetRankedStats(s Summoner, cacheRequests chan CacheRequest) []byte {
	var buffer bytes.Buffer

	stats := getSummonerRankedStats(s.Id)
	leagues := getSummonerLeagues(s.Id)
	soloTierDiv := "Cardboard V"
	teamTier := "Cardboard"
	teamDivision := "V"
	for _, league := range leagues[strconv.FormatInt(s.Id, 10)] {
		if league.Queue == "RANKED_SOLO_5x5" {
			soloTierDiv = fmt.Sprintf("%s %s", league.Tier, league.Entries[0].Division)
		} else if league.Queue == "RANKED_TEAM_5x5" {
			if CompareRanked(league.Tier, league.Entries[0].Division, teamTier, teamDivision) == 1 {
				teamTier = league.Tier
				teamDivision = league.Entries[0].Division
			}
		}
	}
	buffer.WriteString(fmt.Sprintf("%s:\n  Solo Queue League: %s\n  Best Ranked 5's League: %s %s\nChampion Stats:\n", s.Name, soloTierDiv, teamTier, teamDivision))
	for _, champStats := range stats.Champions {
		if champStats.Id > 0 {
			champ, champErr := GetChampion(champStats.Id, cacheRequests)
			if champErr != nil {
				// Not much I can do here...
			}
			buffer.WriteString(fmt.Sprintf("  Champ: %s,", champ.Name))
		} else {
			buffer.WriteString(fmt.Sprintf("\n  Champion Totals: "))
		}
		buffer.WriteString(fmt.Sprintf(
			" Wins: %d, Losses: %d, Kills: %d, Assists: %d, Deaths: %d\n",
			champStats.Stats.TotalSessionsWon, champStats.Stats.TotalSessionsLost, champStats.Stats.TotalChampionKills, champStats.Stats.TotalAssists, champStats.Stats.TotalDeathsPerSession))
	}

	return buffer.Bytes()
}
