package main

type Summoner struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
}

type RankedStats struct {
	SummonerId int64           `json:"summonerId"`
	Name       string          `json:"name"`
	Champions  []ChampionStats `json:"champions"`
	ModifyDate int64           `json:"modifyDate"`
}

type ChampionStats struct {
	Id    int64           `json:"id"`
	Stats AggregatedStats `json:"stats"`
}

type ChampionInfoList struct {
	Champions []ChampionInfo `json:"Champions"`
}

type ChampionList struct {
	Data map[string]Champion `json:"data"`
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
	Id                int64 `json:"id"`
	Active            bool  `json:"active"`
	BotEnabled        bool  `json:"botEnabled"`
	BotMmEnabled      bool  `json:"botMmEnabled"`
	FreeToPlay        bool  `json:"freeToPlay"`
	RankedPlayEnabled bool  `json:"rankedPlayEnabled"`
}

type AggregatedStats struct {
	AverageAssists              int `json:"averageAssists"`
	AverageChampionsKilled      int `json:"averageChampionsKilled"`
	AverageCombatPlayerScore    int `json:"averageCombatPlayerScore"`
	AverageNodeCapture          int `json:"averageNodeCapture"`
	AverageNodeCaptureAssist    int `json:"averageNodeCaptureAssist"`
	AverageNodeNeutralize       int `json:"averageNodeNeutralize"`
	AverageNodeNeutralizeAssist int `json:"averageNodeNeutralizeAssist"`
	AverageNumDeaths            int `json:"averageNumDeaths"`
	AverageObjectivePlayerScore int `json:"averageObjectivePlayerScore"`
	AverageTeamObjective        int `json:"averageTeamObjective"`
	AverageTotalPlayerScore     int `json:"averageTotalPlayerScore"`
	BotGamesPlayed              int `json:"botGamesPlayed"`
	KillingSpree                int `json:"killingSpree"`
	MaxAssists                  int `json:"maxAssists"`
	MaxChampionsKilled          int `json:"maxChampionsKilled"`
	MaxCombatPlayerScore        int `json:"maxCombatPlayerScore"`
	MaxLargestCriticalStrike    int `json:"maxLargestCriticalStrike"`
	MaxLargestKillingSpree      int `json:"maxLargestKillingSpree"`
	MaxNodeCapture              int `json:"maxNodeCapture"`
	MaxNodeCaptureAssist        int `json:"maxNodeCaptureAssist"`
	MaxNodeNeutralize           int `json:"maxNodeNeutralize"`
	MaxNodeNeutralizeAssist     int `json:"maxNodeNeutralizeAssist"`
	MaxObjectivePlayerScore     int `json:"maxObjectivePlayerScore"`
	MaxTeamObjective            int `json:"maxTeamObjective"`
	MaxTimePlayed               int `json:"maxTimePlayed"`
	MaxTimeSpentLiving          int `json:"maxTimeSpentLiving"`
	MaxTotalPlayerScore         int `json:"maxTotalPlayerScore"`
	MostChampionKillsPerSession int `json:"mostChampionKillsPerSession"`
	MostSpellsCast              int `json:"mostSpellsCast"`
	NormalGamesPlayed           int `json:"normalGamesPlayed"`
	RankedPremadeGamesPlayed    int `json:"rankedPremadeGamesPlayed"`
	RankedSoloGamesPlayed       int `json:"rankedSoloGamesPlayed"`
	TotalAssists                int `json:"totalAssists"`
	TotalChampionKills          int `json:"totalChampionKills"`
	TotalDamageDealt            int `json:"totalDamageDealt"`
	TotalDamageTaken            int `json:"totalDamageTaken"`
	TotalDeathsPerSession       int `json:"totalDeathsPerSession"`
	TotalDoubleKills            int `json:"totalDoubleKills"`
	TotalFirstBlood             int `json:"totalFirstBlood"`
	TotalGoldEarned             int `json:"totalGoldEarned"`
	TotalHeal                   int `json:"totalHeal"`
	TotalMagicDamageDealt       int `json:"totalMagicDamageDealt"`
	TotalMinionKills            int `json:"totalMinionKills"`
	TotalNeutralMinionsKilled   int `json:"totalNeutralMinionsKilled"`
	TotalNodeCapture            int `json:"totalNodeCapture"`
	TotalNodeNeutralize         int `json:"totalNodeNeutralize"`
	TotalPentaKills             int `json:"totalPentaKills"`
	TotalPhysicalDamageDealt    int `json:"totalPhysicalDamageDealt"`
	TotalQuadraKills            int `json:"totalQuadraKills"`
	TotalSessionsLost           int `json:"totalSessionsLost"`
	TotalSessionsPlayed         int `json:"totalSessionsPlayed"`
	TotalSessionsWon            int `json:"totalSessionsWon"`
	TotalTripleKills            int `json:"totalTripleKills"`
	TotalTurretsKilled          int `json:"totalTurretsKilled"`
	TotalUnrealKills            int `json:"totalUnrealKills"`
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
