package LeagueDataCache

import (
	"encoding/json"
	"fmt"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"io/ioutil"
	"os"
	"strconv"
)

func storeChampions(file string) {
	cachedChamps := make(map[string]lapi.Champion, len(allChampions))
	for key, val := range allChampions {
		cachedChamps[fmt.Sprintf("%d", key)] = val
	}
	jsonData, err := json.Marshal(cachedChamps)
	if err != nil {
		fmt.Printf("Failed to serialize: %s\n", err.Error())
		return
	}
	os.Remove(file)
	ioutil.WriteFile(file, jsonData, os.ModePerm)
}

func loadChampions(file string) {
	allChampions = make(map[int64]lapi.Champion, 100)
	champData, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		fmt.Printf("Error loading champions: %s\n", readErr.Error())
		return
	}
	var cachedChamps map[string]lapi.Champion
	marshErr := json.Unmarshal(champData, &cachedChamps)
	for key, val := range cachedChamps {
		newKey, _ := strconv.ParseInt(key, 10, 64)
		allChampions[newKey] = val
	}
	if marshErr != nil {
		fmt.Printf("Error loading champions: %s\n", marshErr.Error())
	}
}

func fetchAndCacheChampion(id int64, api *lapi.LolFetcher) lapi.Champion {
	// TODO: re-fetch champions if there is a new one.
	if allChampions == nil || len(allChampions) == 0 {
		champs := api.GetAllChampions()
		for _, champion := range champs.Data {
			allChampions[champion.Id] = champion
		}
	}
	champ, _ := allChampions[id]
	return champ
}
