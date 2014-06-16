package LeagueDataCache

import (
	"encoding/json"
	"fmt"

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

func (mk MatchKey) String() string {
	return fmt.Sprintf("%d.%d", mk.MatchId, mk.SummonerId)
}

type MatchSimple struct {
	ChampionName  string     // Champion ID associated with game.
	ChampionImage string     // URL to fetch the image of the champion
	CreateDate    int64      // Date that end game data was recorded, specified as epoch milliseconds.
	GameId        int64      // Game ID.
	GameMode      string     // Game mode. (legal values: CLASSIC, ODIN, ARAM, TUTORIAL, ONEFORALL, FIRSTBLOOD)
	GameType      string     // Game type. (legal values: CUSTOM_GAME, MATCHED_GAME, TUTORIAL_GAME)
	Invalid       bool       // Invalid flag.
	IpEarned      int        // IP Earned.
	MapId         int        // Map ID.
	Stats         MatchStats // Important Stats
	SubType       string     // Game sub-type. (legal values: NONE, NORMAL, BOT, RANKED_SOLO_5x5, RANKED_PREMADE_3x3, RANKED_PREMADE_5x5, ODIN_UNRANKED, RANKED_TEAM_3x3, RANKED_TEAM_5x5, NORMAL_3x3, BOT_3x3, CAP_5x5, ARAM_UNRANKED_5x5, ONEFORALL_5x5, FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF, URF_BOT)
	Side          string     // blue or purple
}

type MatchStats struct {
	Assists         int
	ChampionsKilled int
	NumDeaths       int
	Win             bool
}

type MatchDetail struct {
	ChampionName  string        // Champion ID associated with game.
	ChampionImage string        // URL to fetch the image of the champion.
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
	Items         []ItemDetail  // List of items bought by player.
	SubType       string        // Game sub-type. (legal values: NONE, NORMAL, BOT, RANKED_SOLO_5x5, RANKED_PREMADE_3x3, RANKED_PREMADE_5x5, ODIN_UNRANKED, RANKED_TEAM_3x3, RANKED_TEAM_5x5, NORMAL_3x3, BOT_3x3, CAP_5x5, ARAM_UNRANKED_5x5, ONEFORALL_5x5, FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF, URF_BOT)
	Side          string        // blue or purple
}

func (md MatchDetail) toCachedMatch(summonerId int64) (cmd cachedMatchDetail) {
	cmd.SummonerId = summonerId
	cmd.GameId = md.GameId
	cmd.PlayedDate = md.CreateDate
	jData, mErr := json.Marshal(md)
	if mErr != nil {
		return
	}
	cmd.Data = jData
	cmd.CacheExpireDate = getExpireTime(false)
	return
}

func (g MatchDetail) toMatchSimple() (ms MatchSimple) {
	ms.CreateDate = g.CreateDate
	ms.ChampionName = g.ChampionName
	ms.ChampionImage = g.ChampionImage
	ms.GameId = g.GameId
	ms.GameMode = g.GameMode
	ms.GameType = g.GameType
	ms.Invalid = g.Invalid
	ms.IpEarned = g.IpEarned
	ms.MapId = g.MapId
	ms.Side = g.Side
	ms.Stats = MatchStats{Assists: g.Stats.Assists, ChampionsKilled: g.Stats.ChampionsKilled, NumDeaths: g.Stats.NumDeaths, Win: g.Stats.Win}
	ms.SubType = g.SubType
	return
}

type cachedMatchDetail struct {
	SummonerId      int64
	GameId          int64
	PlayedDate      int64
	CacheExpireDate int64
	Data            []byte
}

func (cmd cachedMatchDetail) ToMatchDetail() (md MatchDetail, e error) {
	e = json.Unmarshal(cmd.Data, &md)
	return
}

func (cmd cachedMatchDetail) ToMatchSimple() (ms MatchSimple, e error) {
	g, err := cmd.ToMatchDetail()
	return g.toMatchSimple(), err
}

func (cmd cachedMatchDetail) KeyString() string {
	return fmt.Sprintf("%d.%d", cmd.GameId, cmd.SummonerId)
}

type ItemDetail struct {
	Name     string
	ImageUrl string
}

type Player struct {
	ChampionId    int64  // Champion id associated with player.
	ChampionName  string // Champion name
	ChampionImage string // URL to fetch the image of the champion.
	SummonerId    int64  // Summoner id associated with player.
	SummonerName  string // Summoner name of the player
	Side          string // blue or purple
}

type LFScores struct {
	SummonerId int64
	Tags       map[string]LFScore
}

type LFScore struct {
	Tag   string
	Score float64
}
