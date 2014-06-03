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

type LocalMatchHistory struct {
	SummonerId int64
	Games      []LocalMatchSimple
}
type LocalMatchSimple struct {
	ChampionName string          // Champion ID associated with game.
	CreateDate   int64           // Date that end game data was recorded, specified as epoch milliseconds.
	GameId       int64           // Game ID.
	GameMode     string          // Game mode. (legal values: CLASSIC, ODIN, ARAM, TUTORIAL, ONEFORALL, FIRSTBLOOD)
	GameType     string          // Game type. (legal values: CUSTOM_GAME, MATCHED_GAME, TUTORIAL_GAME)
	Invalid      bool            // Invalid flag.
	IpEarned     int             // IP Earned.
	MapId        int             // Map ID.
	Stats        LocalMatchStats // Important Stats
	SubType      string          // Game sub-type. (legal values: NONE, NORMAL, BOT, RANKED_SOLO_5x5, RANKED_PREMADE_3x3, RANKED_PREMADE_5x5, ODIN_UNRANKED, RANKED_TEAM_3x3, RANKED_TEAM_5x5, NORMAL_3x3, BOT_3x3, CAP_5x5, ARAM_UNRANKED_5x5, ONEFORALL_5x5, FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF, URF_BOT)
	Side         string          // blue or purple
}

func LocalMatchSimpleFromGame(g Game) (lm LocalMatchSimple) {
	lm.CreateDate = g.CreateDate
	lm.GameId = g.GameId
	lm.GameMode = g.GameMode
	lm.GameType = g.GameType
	lm.Invalid = g.Invalid
	lm.IpEarned = g.IpEarned
	lm.MapId = g.MapId
	if g.TeamId == 100 {
		lm.Side = "blue"
	} else {
		lm.Side = "purple"
	}

	return
}

type LocalMatchStats struct {
	Assists         int
	ChampionsKilled int
	NumDeaths       int
	Win             bool
}

type LocalMatchDetail struct {
}
