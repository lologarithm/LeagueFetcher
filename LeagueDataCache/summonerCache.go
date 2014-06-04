package LeagueDataCache

import (
	"encoding/json"
	"fmt"
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
	"io/ioutil"
	"os"
	"time"
)

func fetchAndCacheSummoner(name string) lapi.Summoner {
	if summ, ok := allSummonersByName[name]; !ok {
		// First check memcache if in appengine
		fetchSummoner(name)
	} else {
		// TODO: Make sure this works correctly
		if summ.RevisionDate+3600000 < time.Now().Unix()*1000 {
			fetchSummoner(name)
		}
	}
	summ, _ := allSummonersByName[name]
	return summ
}

func fetchSummoner(name string) {
	summoners := lapi.GetSummonerByName(name)
	if s, gotOk := summoners[name]; gotOk {
		allSummonersByName[name] = s
		allSummonersById[s.Id] = s
	}
}

// Gets a list of summoners by Id from cache or remote server. This function needs to be fixed up to be better.
func fetchSummonersById(ids []int64) []lapi.Summoner {
	missingIds := []int64{}
	summoners := make([]lapi.Summoner, len(ids))
	// First check for any cached summoners
	for index, id := range ids {
		if s, ok := allSummonersById[id]; !ok {
			missingIds = append(missingIds, id)
		} else {
			summoners[index] = s
		}
	}

	// Now fetch the summoners that we didn't find
	if len(missingIds) > 0 {
		fetchedSummoners := lapi.GetSummonersById(missingIds)
		for _, value := range fetchedSummoners {
			allSummonersById[value.Id] = value
			allSummonersByName[value.Name] = value
			// Put this summoner where he belongs. Terrible efficiency.
			for i := 0; i < len(ids); i++ {
				if ids[i] == value.Id {
					summoners[i] = value
					break
				}
			}
		}
	}
	return summoners
}

func storeSummoners(file string) {
	jsonData, err := json.Marshal(allSummonersByName)
	if err != nil {
		fmt.Printf("Failed to serialize: %s\n", err.Error())
		return
	}
	os.Remove(file)
	ioutil.WriteFile(file, jsonData, os.ModePerm)
}

func loadSummoners(file string) {
	allSummonersByName = make(map[string]lapi.Summoner, 1)
	allSummonersById = make(map[int64]lapi.Summoner, 1)
	summonerData, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		fmt.Printf("Error loading summoners: %s\n", readErr.Error())
		return
	}
	marshErr := json.Unmarshal(summonerData, &allSummonersByName)
	if marshErr != nil {
		fmt.Printf("Error loading summoners: %s\n", marshErr.Error())
	}

	for _, summ := range allSummonersByName {
		allSummonersById[summ.Id] = summ
	}
}
