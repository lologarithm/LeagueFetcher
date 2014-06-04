package LeagueApi

type Summoner struct {
	Id            int64
	Name          string
	ProfileIconId int
	SummonerLevel int
	RevisionDate  int64
}

type RankedStats struct {
	SummonerId int64
	Name       string
	Champions  []ChampionStats
	ModifyDate int64
}

type ChampionStats struct {
	Id           int64
	ChampionName string
	Stats        AggregatedStats
}

type ChampionInfoList struct {
	Champions []ChampionInfo
}

type ChampionList struct {
	Data map[string]Champion
}

type Champion struct {
	Id    int64
	Title string

	Allytips  []string
	Blurb     string
	Enemytips []string
	//image	ImageDto
	//info	InfoDto
	Key     string
	Lore    string
	Name    string
	Partype string
	//passive	PassiveDto
	//recommended	List[RecommendedDto]
	//skins	List[SkinDto]
	Spells []ChampionSpell
	//stats	StatsDto
	Tags []string
}

type ChampionSpell struct {
	Key  string
	Name string

	//altimages	List[ImageDto]
	Cooldown     []float64
	CooldownBurn string
	Cost         []int
	CostBurn     string
	CostType     string
	Description  string
	//effect	List[object]	This field is a List of List of Integer.
	EffectBurn []string
	//image	ImageDto
	//leveltip	LevelTipDto
	Maxrank int
	//range	object	This field is either a List of Integer or the String 'self' for spells that target one's own champion.
	RangeBurn            string
	Resource             string
	SanitizedDescription string
	SanitizedTooltip     string
	Tooltip              string
	//vars	List[SpellVarsDto]
}

type ChampionInfo struct {
	Id                int64
	Active            bool
	BotEnabled        bool
	BotMmEnabled      bool
	FreeToPlay        bool
	RankedPlayEnabled bool
}

type AggregatedStats struct {
	AverageAssists              int
	AverageChampionsKilled      int
	AverageCombatPlayerScore    int
	AverageNodeCapture          int
	AverageNodeCaptureAssist    int
	AverageNodeNeutralize       int
	AverageNodeNeutralizeAssist int
	AverageNumDeaths            int
	AverageObjectivePlayerScore int
	AverageTeamObjective        int
	AverageTotalPlayerScore     int
	BotGamesPlayed              int
	KillingSpree                int
	MaxAssists                  int
	MaxChampionsKilled          int
	MaxCombatPlayerScore        int
	MaxLargestCriticalStrike    int
	MaxLargestKillingSpree      int
	MaxNodeCapture              int
	MaxNodeCaptureAssist        int
	MaxNodeNeutralize           int
	MaxNodeNeutralizeAssist     int
	MaxObjectivePlayerScore     int
	MaxTeamObjective            int
	MaxTimePlayed               int
	MaxTimeSpentLiving          int
	MaxTotalPlayerScore         int
	MostChampionKillsPerSession int
	MostSpellsCast              int
	NormalGamesPlayed           int
	RankedPremadeGamesPlayed    int
	RankedSoloGamesPlayed       int
	TotalAssists                int
	TotalChampionKills          int
	TotalDamageDealt            int
	TotalDamageTaken            int
	TotalDeathsPerSession       int
	TotalDoubleKills            int
	TotalFirstBlood             int
	TotalGoldEarned             int
	TotalHeal                   int
	TotalMagicDamageDealt       int
	TotalMinionKills            int
	TotalNeutralMinionsKilled   int
	TotalNodeCapture            int
	TotalNodeNeutralize         int
	TotalPentaKills             int
	TotalPhysicalDamageDealt    int
	TotalQuadraKills            int
	TotalSessionsLost           int
	TotalSessionsPlayed         int
	TotalSessionsWon            int
	TotalTripleKills            int
	TotalTurretsKilled          int
	TotalUnrealKills            int
}

type League struct {
	Entries       []LeagueEntry // The requested league entries.
	Name          string        // This name is an internal place-holder name only. Display and localization of names in the game client are handled client-side.
	ParticipantId string        // Specifies the relevant participant that is a member of this league (i.e., a requested summoner ID, a requested team ID, or the ID of a team to which one of the requested summoners belongs). Only present when full league is requested so that participant's entry can be identified. Not present when individual entry is requested.
	Queue         string        // The league's queue type. (legal values: RANKED_SOLO_5x5, RANKED_TEAM_3x3, RANKED_TEAM_5x5)
	Tier          string        // The league's tier. (legal values: CHALLENGER, DIAMOND, PLATINUM, GOLD, SILVER, BRONZE)
}

type LeagueEntry struct {
	Division         string     // The league division of the participant.
	IsFreshBlood     bool       // Specifies if the participant is fresh blood.
	IsHotStreak      bool       // Specifies if the participant is on a hot streak.
	IsInactive       bool       // Specifies if the participant is inactive.
	IsVeteran        bool       // Specifies if the participant is a veteran.
	LeaguePoints     int        // The league points of the participant.
	MiniSeries       MiniSeries // Mini series data for the participant. Only present if the participant is currently in a mini series.
	PlayerOrTeamId   string     // The ID of the participant (i.e., summoner or team) represented by this entry.
	PlayerOrTeamName string     // The name of the the participant (i.e., summoner or team) represented by this entry.
	Wins             int        // The number of wins for the participant.
}

type MiniSeries struct {
	Losses   int    // Number of current losses in the mini series.
	Progress string // String showing the current, sequential mini series progress where 'W' represents a win, 'L' represents a loss, and 'N' represents a game that hasn't been played yet.
	Target   int    // Number of wins required for promotion.
	Wins     int    // Number of current wins in the mini series.
}

type Team struct {
	CreateDate                    int64
	FullId                        string
	LastGameDate                  int64
	LastJoinDate                  int64
	LastJoinedRankedTeamQueueDate int64
	MatchHistory                  []MatchHistorySummary
	ModifyDate                    int64
	Name                          string
	Roster                        Roster
	SecondLastJoinDate            int64
	Status                        string
	Tag                           string
	TeamStatDetails               []TeamStatDetail
	ThirdLastJoinDate             int64
}

type Roster struct {
	MemberList []TeamMemberInfo
	OwnerId    int64
}

type TeamMemberInfo struct {
	inviteDate int64
	joinDate   int64
	playerId   int64
	status     string
}

type MatchHistorySummary struct {
	Assists           int
	Date              int64
	Deaths            int
	GameId            int64
	GameMode          string
	Invalid           bool
	Kills             int
	MapId             int
	OpposingTeamKills int
	OpposingTeamName  string
	Win               bool
}

type TeamStatDetail struct {
	AverageGamesPlayed int
	Losses             int
	TeamStatType       string
	Wins               int
}

type PlayerStatsSummaryList struct {
	PlayerStatSummaries []PlayerStatsSummary // Collection of player stats summaries associated with the summoner.
	SummonerId          int64                // Summoner ID.
}

type PlayerStatsSummary struct {
	AggregatedStats       AggregatedStats // Aggregated stats.
	Losses                int             // Number of losses for this queue type. Returned for ranked queue types only.
	ModifyDate            int64           // Date stats were last modified specified as epoch milliseconds.
	PlayerStatSummaryType string          // Player stats summary type. (legal values: AramUnranked5x5, CoopVsAI, CoopVsAI3x3, OdinUnranked, RankedPremade3x3, RankedPremade5x5, RankedSolo5x5, RankedTeam3x3, RankedTeam5x5, Unranked, Unranked3x3, OneForAll5x5, FirstBlood1x1, FirstBlood2x2, SummonersRift6x6, CAP5x5, URF, URFBots)
	Wins                  int             // Number of wins for this queue type.
}

type RecentGames struct {
	Games      []Game // Collection of recent games played (max 10).
	SummonerId int64  // Summoner ID
}

type Game struct {
	ChampionId    int64    // Champion ID associated with game.
	CreateDate    int64    // Date that end game data was recorded, specified as epoch milliseconds.
	FellowPlayers []Player // Other players associated with the game.
	GameId        int64    // Game ID.
	GameMode      string   // Game mode. (legal values: CLASSIC, ODIN, ARAM, TUTORIAL, ONEFORALL, FIRSTBLOOD)
	GameType      string   // Game type. (legal values: CUSTOM_GAME, MATCHED_GAME, TUTORIAL_GAME)
	Invalid       bool     // Invalid flag.
	IpEarned      int      // IP Earned.
	Level         int      // Level.
	MapId         int      // Map ID.
	Spell1        int      // ID of first summoner spell.
	Spell2        int      // ID of second summoner spell.
	Stats         RawStats // Statistics associated with the game for this summoner.
	SubType       string   // Game sub-type. (legal values: NONE, NORMAL, BOT, RANKED_SOLO_5x5, RANKED_PREMADE_3x3, RANKED_PREMADE_5x5, ODIN_UNRANKED, RANKED_TEAM_3x3, RANKED_TEAM_5x5, NORMAL_3x3, BOT_3x3, CAP_5x5, ARAM_UNRANKED_5x5, ONEFORALL_5x5, FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF, URF_BOT)
	TeamId        int      // Team ID associated with game. Team ID 100 is blue team. Team ID 200 is purple team.
	ExpireTime    int64    // Time this data will expire.
}

type Player struct {
	ChampionId int64 // Champion id associated with player.
	SummonerId int64 // Summoner id associated with player.
	TeamId     int   // Team id associated with player.
}

type RawStats struct {
	Assists                         int
	BarracksKilled                  int // Number of enemy inhibitors killed.
	ChampionsKilled                 int
	CombatPlayerScore               int
	ConsumablesPurchased            int
	DamageDealtPlayer               int
	DoubleKills                     int
	FirstBlood                      int
	Gold                            int
	GoldEarned                      int
	GoldSpent                       int
	Item0                           int
	Item1                           int
	Item2                           int
	Item3                           int
	Item4                           int
	Item5                           int
	Item6                           int
	ItemsPurchased                  int
	KillingSprees                   int
	LargestCriticalStrike           int
	LargestKillingSpree             int
	LargestMultiKill                int
	LegendaryItemsCreated           int // Number of tier 3 items built.
	Level                           int
	MagicDamageDealtPlayer          int
	MagicDamageDealtToChampions     int
	MagicDamageTaken                int
	MinionsDenied                   int
	MinionsKilled                   int
	NeutralMinionsKilled            int
	NeutralMinionsKilledEnemyJungle int
	NeutralMinionsKilledYourJungle  int
	NexusKilled                     bool // Flag specifying if the summoner got the killing blow on the nexus.
	NodeCapture                     int
	NodeCaptureAssist               int
	NodeNeutralize                  int
	NodeNeutralizeAssist            int
	NumDeaths                       int
	NumItemsBought                  int
	ObjectivePlayerScore            int
	PentaKills                      int
	PhysicalDamageDealtPlayer       int
	PhysicalDamageDealtToChampions  int
	PhysicalDamageTaken             int
	QuadraKills                     int
	SightWardsBought                int
	Spell1Cast                      int // Number of times first champion spell was cast.
	Spell2Cast                      int // Number of times second champion spell was cast.
	Spell3Cast                      int // Number of times third champion spell was cast.
	Spell4Cast                      int // Number of times fourth champion spell was cast.
	SummonSpell1Cast                int
	SummonSpell2Cast                int
	SuperMonsterKilled              int
	Team                            int
	TeamObjective                   int
	TimePlayed                      int
	TotalDamageDealt                int
	TotalDamageDealtToChampions     int
	TotalDamageTaken                int
	TotalHeal                       int
	TotalPlayerScore                int
	TotalScoreRank                  int
	TotalTimeCrowdControlDealt      int
	TotalUnitsHealed                int
	TripleKills                     int
	TrueDamageDealtPlayer           int
	TrueDamageDealtToChampions      int
	TrueDamageTaken                 int
	TurretsKilled                   int
	UnrealKills                     int
	VictoryPointTotal               int
	VisionWardsBought               int
	WardKilled                      int
	WardPlaced                      int
	Win                             bool //Flag specifying whether or not this game was won.
}
