package main

import (
	"fmt"
)

func main() {
	SetupCache()

	champ := GetChampion(5)
	fmt.Printf("Champ: %v\n", champ)

	summoner := GetSummoner("lologarithm")
	summoner2 := GetSummoner("comradetaters")
	fmt.Printf("Summoner: %v\n", summoner)
	fmt.Printf("Summoner2: %v\n", summoner2)

	stats := getSummonerRankedStats(summoner.Id)
	fmt.Printf("Stats: %v\n", stats)

	SaveCache()
}
