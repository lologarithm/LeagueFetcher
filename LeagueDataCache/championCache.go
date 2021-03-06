package LeagueDataCache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
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
	if allChampions != nil {
		return
	}
	if _, err := os.Stat(file); err == nil {
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
}

func loadItems(file string) {
	if allItems != nil {
		return
	}
	if _, err := os.Stat(file); err == nil {
		iData, readErr := ioutil.ReadFile(file)
		if readErr != nil {
			fmt.Printf("Error loading items: %s\n", readErr.Error())
			return
		}
		marshErr := json.Unmarshal(iData, &allItems)
		if marshErr != nil {
			fmt.Printf("Error loading items: %s\n", marshErr.Error())
			return
		}
	} else {
		return
	}
	allItems.ItemsById = make(map[int64]lapi.Item, len(allItems.Data))
	if allItems.ItemsById != nil {
		for _, item := range allItems.Data {
			allItems.ItemsById[item.Id] = item
		}
	}
}

func fetchAndCacheChampion(id int64, api *lapi.LolFetcher) (lapi.Champion, error) {
	if allChampions == nil || len(allChampions) == 0 {
		allChampions = make(map[int64]lapi.Champion, 100)
		champs, fetchErr := api.GetAllChampions()
		if fetchErr != nil {
			return lapi.Champion{}, fetchErr
		}
		for _, champion := range champs.Data {
			allChampions[champion.Id] = champion
		}
	}
	champ, ok := allChampions[id]
	if !ok {
		fetchedCh, err := api.GetChampion(id)
		if err != nil {
			return fetchedCh, err
		}
		allChampions[champ.Id] = fetchedCh
		champ = fetchedCh
	}
	return champ, nil
}
