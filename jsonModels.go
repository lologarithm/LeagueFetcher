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
