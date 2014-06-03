package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	SetupCache()

	summoner := GetSummoner("lologarithm")
	summoner2 := GetSummoner("comradetaters")
	fmt.Printf("Summoner: %v\n", summoner)
	fmt.Printf("Summoner2: %v\n", summoner2)

	PrintRankedStats(summoner2)
	PrintSummonerSummary(summoner)
	SaveCache()
}

func PrintSummonerSummary(s Summoner) {
	summary := getSummonerSummaryStats(s.Id)
	for _, stats := range summary.PlayerStatSummaries {
		if strings.Contains(stats.PlayerStatSummaryType, "Ranked") {
			fmt.Printf("Game Type: %s\n  Wins: %d\n  Losses: %d\n  Kills: %d\n  Assists: %d\n", stats.PlayerStatSummaryType, stats.Wins, stats.Losses, stats.AggregatedStats.TotalChampionKills, stats.AggregatedStats.TotalAssists)
		} else {
			fmt.Printf("Game Type: %s\n  Wins: %d\n  Kills: %d\n  Assists: %d\n", stats.PlayerStatSummaryType, stats.Wins, stats.AggregatedStats.TotalChampionKills, stats.AggregatedStats.TotalAssists)
		}
	}
}

func PrintRankedStats(s Summoner) {
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
	fmt.Printf("%s:\n  Solo Queue League: %s\n  Best Ranked 5's League: %s %s\nChampion Stats:\n", s.Name, soloTierDiv, teamTier, teamDivision)
	for _, champStats := range stats.Champions {
		champ := GetChampion(champStats.Id)
		if champ.Id > 0 {
			fmt.Printf("  Champ: %s,", champ.Name)
		} else {
			fmt.Printf("\n  Champion Totals: ")
		}
		fmt.Printf(
			" Wins: %d, Losses: %d, Kills: %d, Assists: %d, Deaths: %d\n",
			champStats.Stats.TotalSessionsWon, champStats.Stats.TotalSessionsLost, champStats.Stats.TotalChampionKills, champStats.Stats.TotalAssists, champStats.Stats.TotalDeathsPerSession)
	}
}
