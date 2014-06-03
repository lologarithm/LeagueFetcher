package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	TierChallenger = 10
	TierDiamond    = 5
	TierPlatinum   = 4
	TierGold       = 3
	TierSilver     = 2
	TierBronze     = 1
)

const (
	baseUrl         = "https://prod.api.pvp.net/api/lol"
	region          = "na"
	summonerVersion = "v1.4"
	statsVersion    = "v1.3"
	champVersion    = "v1.2"
	leagueVersion   = "v2.4"
	teamVersion     = "v2.3"
)

var apiKey string

func makeUrl(version string, method string) string {
	url := fmt.Sprintf("%s/%s/%s/%s?api_key=%s", baseUrl, region, version, method, apiKey)
	//fmt.Printf("URL: %s\n", url)
	return url
}

func makeStaticDataUrl(version string, method string, params string) string {
	url := fmt.Sprintf("%s/static-data/%s/%s/%s?api_key=%s%s", baseUrl, region, version, method, apiKey, params)
	//fmt.Printf("URL: %s\n", url)
	return url
}

func makeRequest(url string, value interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to open conn: %s\n", err.Error())
		return
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		fmt.Printf("Failed to open conn: %s\n", err.Error())
		return
	}

	unmarshErr := json.Unmarshal(body, value)
	if unmarshErr != nil {
		fmt.Printf("Failed to unmarshal json: %s\n", unmarshErr.Error())
	}
}

func getSummonerByName(name string) (summoners map[string]Summoner) {
	makeRequest(makeUrl(summonerVersion, "summoner/by-name/"+name), &summoners)
	return
}

func getSummonerRankedStats(id int64) (srs RankedStats) {
	method := fmt.Sprintf("stats/by-summoner/%d/ranked", id)
	makeRequest(makeUrl(statsVersion, method), &srs)
	return
}

func getSummonerSummaryStats(id int64) (stats PlayerStatsSummaryList) {
	method := fmt.Sprintf("stats/by-summoner/%d/summary", id)
	makeRequest(makeUrl(statsVersion, method), &stats)
	return
}

func getSummonerLeagues(id int64) (leagues map[string][]League) {
	method := fmt.Sprintf("league/by-summoner/%d/entry", id)
	makeRequest(makeUrl(leagueVersion, method), &leagues)
	return
}

func getSummonerTeams(id int64) (teams map[string][]Team) {
	method := fmt.Sprintf("team/by-summoner/%d", id)
	makeRequest(makeUrl(teamVersion, method), &teams)
	return
}

func getAllChampions() (champs ChampionList) {
	params := "&champData=all"
	makeRequest(makeStaticDataUrl(champVersion, "champion", params), &champs)
	return
}

// -1 means tier 1 is worse, 0 means equal, 1 means tier 1 is better
func CompareRanked(tier1 string, div1 string, tier2 string, div2 string) int {
	tierMap := map[string]int{"BRONZE": TierBronze, "SILVER": TierSilver, "GOLD": TierGold, "PLATINUM": TierPlatinum, "DIAMOND": TierDiamond, "CHALLENGER": TierChallenger}
	if tierMap[tier1] > tierMap[tier2] {
		return 1
	} else if tierMap[tier1] < tierMap[tier2] {
		return -1
	}

	if LeagueDivisionToNumber(div1) < LeagueDivisionToNumber(div2) {
		return 1
	} else if LeagueDivisionToNumber(div1) > LeagueDivisionToNumber(div2) {
		return -1
	}

	return 0
}

// Not enough different divisions to need a real roman numeral translator.
func LeagueDivisionToNumber(div string) int {
	switch div {
	case "V":
		return 5
	case "IV":
		return 4
	case "III":
		return 3
	case "II":
		return 2
	case "I":
		return 1
	}
	// Return some very large value
	return 100
}
