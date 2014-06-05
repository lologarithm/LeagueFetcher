package LeagueDataCache

import (
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

// Combines data from Summoner, RankedStats, Solo Queue League, and Ranked Teams
type SummonerRankedData struct {
	// Masteries
	// Runes
	lapi.Summoner
	lapi.RankedStats
	Solo5sLeague      lapi.League
	Solo3sLeague      lapi.League
	RankedTeamLeagues map[string]lapi.League
	ExpireTime        int64 // Unix time of when this data is expired
}

// Combines Team information with the league information
type Team struct {
	lapi.Team
	RankedLeague lapi.League
}

type MatchHistory struct {
	SummonerId int64
	Games      []MatchSimple
}

type MatchKey struct {
	MatchId    int64
	SummonerId int64
}
type MatchSimple struct {
	ChampionName string     // Champion ID associated with game.
	CreateDate   int64      // Date that end game data was recorded, specified as epoch milliseconds.
	GameId       int64      // Game ID.
	GameMode     string     // Game mode. (legal values: CLASSIC, ODIN, ARAM, TUTORIAL, ONEFORALL, FIRSTBLOOD)
	GameType     string     // Game type. (legal values: CUSTOM_GAME, MATCHED_GAME, TUTORIAL_GAME)
	Invalid      bool       // Invalid flag.
	IpEarned     int        // IP Earned.
	MapId        int        // Map ID.
	Stats        MatchStats // Important Stats
	SubType      string     // Game sub-type. (legal values: NONE, NORMAL, BOT, RANKED_SOLO_5x5, RANKED_PREMADE_3x3, RANKED_PREMADE_5x5, ODIN_UNRANKED, RANKED_TEAM_3x3, RANKED_TEAM_5x5, NORMAL_3x3, BOT_3x3, CAP_5x5, ARAM_UNRANKED_5x5, ONEFORALL_5x5, FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF, URF_BOT)
	Side         string     // blue or purple
}

func NewMatchSimpleFromGame(g lapi.Game) (lm MatchSimple) {
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
	lm.Stats = MatchStats{Assists: g.Stats.Assists, ChampionsKilled: g.Stats.ChampionsKilled, NumDeaths: g.Stats.NumDeaths, Win: g.Stats.Win}
	lm.SubType = g.SubType
	return
}

func NewMatchDetailsFromGame(g lapi.Game) (lmd MatchDetail) {
	lmd.CreateDate = g.CreateDate
	lmd.GameId = g.GameId
	lmd.GameMode = g.GameMode
	lmd.GameType = g.GameType
	lmd.Invalid = g.Invalid
	lmd.IpEarned = g.IpEarned
	lmd.MapId = g.MapId
	if g.TeamId == 100 {
		lmd.Side = "blue"
	} else {
		lmd.Side = "purple"
	}
	lmd.SubType = g.SubType
	players := []Player{}
	for _, player := range g.FellowPlayers {
		p := Player{ChampionId: player.ChampionId, SummonerId: player.SummonerId}
		if player.TeamId == 100 {
			p.Side = "blue"
		} else {
			p.Side = "purple"
		}
		players = append(players, p)
	}
	lmd.FellowPlayers = players
	lmd.Spell1 = g.Spell1
	lmd.Spell2 = g.Spell2
	lmd.Stats = g.Stats
	return
}

type MatchStats struct {
	Assists         int
	ChampionsKilled int
	NumDeaths       int
	Win             bool
}

type MatchDetail struct {
	ChampionName  string        // Champion ID associated with game.
	CreateDate    int64         // Date that end game data was recorded, specified as epoch milliseconds.
	FellowPlayers []Player      // Other players associated with the game.
	GameId        int64         // Game ID.
	GameMode      string        // Game mode. (legal values: CLASSIC, ODIN, ARAM, TUTORIAL, ONEFORALL, FIRSTBLOOD)
	GameType      string        // Game type. (legal values: CUSTOM_GAME, MATCHED_GAME, TUTORIAL_GAME)
	Invalid       bool          // Invalid flag.
	IpEarned      int           // IP Earned.
	MapId         int           // Map ID.
	Spell1        int           // ID of first summoner spell.
	Spell2        int           // ID of second summoner spell.
	Stats         lapi.RawStats // Statistics associated with the game for this summoner.
	SubType       string        // Game sub-type. (legal values: NONE, NORMAL, BOT, RANKED_SOLO_5x5, RANKED_PREMADE_3x3, RANKED_PREMADE_5x5, ODIN_UNRANKED, RANKED_TEAM_3x3, RANKED_TEAM_5x5, NORMAL_3x3, BOT_3x3, CAP_5x5, ARAM_UNRANKED_5x5, ONEFORALL_5x5, FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF, URF_BOT)
	Side          string        // blue or purple
}

type Player struct {
	ChampionName string // Champion name
	ChampionId   int64  // Champion id associated with player.
	SummonerId   int64  // Summoner id associated with player.
	SummonerName string // Summoner name of the player
	Side         string // blue or purple
}
