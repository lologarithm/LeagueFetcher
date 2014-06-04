package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	championCache = "ccache.json"
	summonerCache = "scache.json"
)

var allSummoners map[string]Summoner
var allChampions map[int64]Champion
var allLeagues map[int64]League
var allTeams map[int64]Team

func storeSummoners(file string) {
	jsonData, err := json.Marshal(allSummoners)
	if err != nil {
		fmt.Printf("Failed to serialize: %s\n", err.Error())
		return
	}
	os.Remove(file)
	ioutil.WriteFile(file, jsonData, os.ModePerm)
}

func loadSummoners(file string) {
	allSummoners = make(map[string]Summoner, 2)
	summonerData, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		fmt.Printf("Error loading summoners: %s\n", readErr.Error())
		return
	}
	marshErr := json.Unmarshal(summonerData, &allSummoners)
	if marshErr != nil {
		fmt.Printf("Error loading summoners: %s\n", marshErr.Error())
	}
}

func storeChampions(file string) {
	cachedChamps := make(map[string]Champion, len(allChampions))
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
	allChampions = make(map[int64]Champion, 100)
	champData, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		fmt.Printf("Error loading champions: %s\n", readErr.Error())
		return
	}
	var cachedChamps map[string]Champion
	marshErr := json.Unmarshal(champData, &cachedChamps)
	for key, val := range cachedChamps {
		newKey, _ := strconv.ParseInt(key, 10, 64)
		allChampions[newKey] = val
	}
	if marshErr != nil {
		fmt.Printf("Error loading champions: %s\n", marshErr.Error())
	}
}

func GetSummoner(name string) Summoner {
	if _, ok := allSummoners[name]; !ok {
		summoners := getSummonerByName(name)
		if s, gotOk := summoners[name]; gotOk {
			allSummoners[name] = s
		}
	}
	summ, _ := allSummoners[name]
	return summ
}

func GetChampion(id int64) Champion {
	if allChampions == nil || len(allChampions) == 0 {
		champs := getAllChampions()
		for _, champion := range champs.Data {
			// fmt.Printf("Got Champ %s: %v\n", key, champion)
			allChampions[champion.Id] = champion
		}
	}
	champ, _ := allChampions[id]
	return champ
}

func SetupCache() {
	loadChampions(championCache)
	loadSummoners(summonerCache)
}

func SaveCache() {
	storeChampions(championCache)
	storeSummoners(summonerCache)
}
