package main

// Combines data from Summoner, RankedStats, Solo Queue League, and Ranked Teams
type LocalSummoner struct {
	// Masteries
	// Runes
	Summoner
	RankedStats
	SoloLeague League
	Teams      map[string]LocalTeam
}

// Combines Team information with the league information
type LocalTeam struct {
	Team
	RankedLeague League
}
