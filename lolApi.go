package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseUrl         = "https://prod.api.pvp.net/api/lol"
	region          = "na"
	summonerVersion = "v1.4"
	statsVersion    = "v1.3"
	champVersion    = "v1.2"
	apiKey          = "SETME"
)

func makeUrl(version string, method string) string {
	url := fmt.Sprintf("%s/%s/%s/%s?api_key=%s", baseUrl, region, version, method, apiKey)
	fmt.Printf("URL: %s\n", url)
	return url
}

func makeStaticDataUrl(version string, method string, params string) string {
	url := fmt.Sprintf("%s/static-data/%s/%s/%s?api_key=%s%s", baseUrl, region, version, method, apiKey, params)
	fmt.Printf("URL: %s\n", url)
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

func getAllChampions() (champs ChampionList) {
	params := "&champData=all"
	makeRequest(makeStaticDataUrl(champVersion, "champion", params), &champs)
	return
}
