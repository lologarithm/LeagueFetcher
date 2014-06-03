package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const (
	championCache = "ccache.json"
	summonerCache = "scache.json"
)

type CacheRequest struct {
	Response chan CacheResponse
	Type     string
	Key      string
}

type CacheResponse struct {
	Value interface{}
	Ok    bool
}

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

func fetchAndCacheSummoner(name string) Summoner {

	if summ, ok := allSummoners[name]; !ok {
		// First check memcache if in appengine
		fetchSummoner(name)
	} else {
		// TODO: Make sure this works correctly
		if summ.RevisionDate+3600000 < time.Now().Unix()*1000 {
			fetchSummoner(name)
		}
	}
	summ, _ := allSummoners[name]
	return summ
}

func fetchSummoner(name string) {
	summoners := getSummonerByName(name)
	if s, gotOk := summoners[name]; gotOk {
		allSummoners[name] = s
	}
}
func fetchAndCacheChampion(id int64) Champion {
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

func RunCache(exit chan bool, requests chan CacheRequest) {
	setupCache()
	for {
		select {
		case <-exit:
			return
		case cRequest := <-requests:
			fetchCache(cRequest)
		}
	}
	saveCache()
}

func fetchCache(request CacheRequest) {
	response := CacheResponse{Ok: false}
	switch request.Type {
	case "summoner":
		summoner := fetchAndCacheSummoner(request.Key)
		response.Ok = true
		response.Value = summoner
	case "champion":
		intKey, err := strconv.ParseInt(request.Key, 10, 64)
		if err != nil {
			break
		}
		champ := fetchAndCacheChampion(intKey)
		response.Value = champ
		response.Ok = true
	}

	request.Response <- response
}

func setupCache() {
	loadChampions(championCache)
	loadSummoners(summonerCache)
}

func saveCache() {
	storeChampions(championCache)
	storeSummoners(summonerCache)
}

type cacheError struct {
	message string
}

func (c cacheError) Error() string {
	return c.message
}

// Public Functions that use channel to communicate with cache goroutine

// Fetch Champion from cache goroutine.
func GetChampion(id int64, cacheRequests chan CacheRequest) (c Champion, e error) {
	result := make(chan CacheResponse, 1)
	cReq := CacheRequest{Type: "champion", Key: fmt.Sprintf("%d", id), Response: result}
	cacheRequests <- cReq
	champResponse := <-result
	if champResponse.Ok {
		champ, _ := champResponse.Value.(Champion)
		return champ, nil
	}

	return Champion{}, cacheError{message: "Failed to retrieve champion."}
}

// Fetch Summoner from cache goroutine
func GetSummoner(name string, cacheRequests chan CacheRequest) (s Summoner, e error) {
	summResponse := make(chan CacheResponse, 1)
	summRequest := CacheRequest{Response: summResponse, Type: "summoner", Key: name}
	cacheRequests <- summRequest
	respValue := <-summResponse
	if respValue.Ok {
		summoner, _ := respValue.Value.(Summoner)
		return summoner, nil
	}
	return Summoner{}, cacheError{message: "Failed to retrieve summoner."}
}
