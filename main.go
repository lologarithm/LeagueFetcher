package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseUrl         = "https://prod.api.pvp.net/api/lol/"
	region          = "na"
	summonerVersion = "v1.4"
	statsVersion    = "v1.3"
	apiKey          = "?api_key=f0fa5a3c-e718-4004-b264-b9a64fc7a444"
)

func main() {

}

func getUrl(method string) string {
	url := baseUrl + region + "/" + version + "/" + method + apiKey
	fmt.Printf("URL: %s\n", url)
	return url
}

func makeRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to open conn: %s\n", err.Error())
		return []byte{}
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		fmt.Printf("Failed to open conn: %s\n", err.Error())
		return []byte{}
	}
	return body
}

func getSummonerByName(name string) (summoners map[string]Summoner) {
	response := makeRequest(getUrl("summoner/by-name/" + name))

	if len(response) == 0 {
		return
	}

	unmarshErr := json.Unmarshal(response, &summoners)
	if unmarshErr != nil {
		fmt.Printf("Failed to unmarshal json: %s\n", unmarshErr.Error())
	}

	return
}

func getSummonerRankedStats(id int) {

}
