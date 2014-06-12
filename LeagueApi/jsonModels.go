package LeagueApi

import (
	"fmt"
)

type Summoner struct {
	Id             int64
	Name           string
	NormalizedName string
	ProfileIconId  int   `datastore:",noindex"`
	SummonerLevel  int   `datastore:",noindex"`
	RevisionDate   int64 `datastore:",noindex"`
}

type RankedStats struct {
	SummonerId int64
	Name       string
	Champions  []ChampionStats
	ModifyDate int64
}

type ChampionStats struct {
	Id            int64
	ChampionName  string
	ChampionImage string
	Stats         AggregatedStats
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

	Allytips    []string
	Blurb       string
	Enemytips   []string
	Image       Image
	Info        Info
	Key         string
	Lore        string
	Name        string
	Partype     string
	Passive     Passive
	Recommended []Recommended
	Skins       []Skin
	Spells      []ChampionSpell
	Stats       Stats
	Tags        []string
}

type ChampionSpell struct {
	Key                  string
	Name                 string
	Altimages            []Image
	Cooldown             []float64
	CooldownBurn         string
	Cost                 []int
	CostBurn             string
	CostType             string
	Description          string
	Effect               [][]int
	EffectBurn           []string
	Image                Image
	Leveltip             LevelTip
	Maxrank              int
	Range                interface{} //	This field is either a List of Integer or the String 'self' for spells that target one's own champion.
	RangeBurn            string
	Resource             string
	SanitizedDescription string
	SanitizedTooltip     string
	Tooltip              string
	Vars                 []SpellVars
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
	AverageAssists              int `json:",omitempty"`
	AverageChampionsKilled      int `json:",omitempty"`
	AverageCombatPlayerScore    int `json:",omitempty"`
	AverageNodeCapture          int `json:",omitempty"`
	AverageNodeCaptureAssist    int `json:",omitempty"`
	AverageNodeNeutralize       int `json:",omitempty"`
	AverageNodeNeutralizeAssist int `json:",omitempty"`
	AverageNumDeaths            int `json:",omitempty"`
	AverageObjectivePlayerScore int `json:",omitempty"`
	AverageTeamObjective        int `json:",omitempty"`
	AverageTotalPlayerScore     int `json:",omitempty"`
	BotGamesPlayed              int `json:",omitempty"`
	KillingSpree                int `json:",omitempty"`
	MaxAssists                  int `json:",omitempty"`
	MaxChampionsKilled          int `json:",omitempty"`
	MaxCombatPlayerScore        int `json:",omitempty"`
	MaxLargestCriticalStrike    int `json:",omitempty"`
	MaxLargestKillingSpree      int `json:",omitempty"`
	MaxNodeCapture              int `json:",omitempty"`
	MaxNodeCaptureAssist        int `json:",omitempty"`
	MaxNodeNeutralize           int `json:",omitempty"`
	MaxNodeNeutralizeAssist     int `json:",omitempty"`
	MaxObjectivePlayerScore     int `json:",omitempty"`
	MaxTeamObjective            int `json:",omitempty"`
	MaxTimePlayed               int `json:",omitempty"`
	MaxTimeSpentLiving          int `json:",omitempty"`
	MaxTotalPlayerScore         int `json:",omitempty"`
	MostChampionKillsPerSession int `json:",omitempty"`
	MostSpellsCast              int `json:",omitempty"`
	NormalGamesPlayed           int `json:",omitempty"`
	RankedPremadeGamesPlayed    int `json:",omitempty"`
	RankedSoloGamesPlayed       int `json:",omitempty"`
	TotalAssists                int
	TotalChampionKills          int
	TotalDamageDealt            int `json:",omitempty"`
	TotalDamageTaken            int `json:",omitempty"`
	TotalDeathsPerSession       int
	TotalDoubleKills            int `json:",omitempty"`
	TotalFirstBlood             int `json:",omitempty"`
	TotalGoldEarned             int `json:",omitempty"`
	TotalHeal                   int `json:",omitempty"`
	TotalMagicDamageDealt       int `json:",omitempty"`
	TotalMinionKills            int `json:",omitempty"`
	TotalNeutralMinionsKilled   int `json:",omitempty"`
	TotalNodeCapture            int `json:",omitempty"`
	TotalNodeNeutralize         int `json:",omitempty"`
	TotalPentaKills             int `json:",omitempty"`
	TotalPhysicalDamageDealt    int `json:",omitempty"`
	TotalQuadraKills            int `json:",omitempty"`
	TotalSessionsLost           int
	TotalSessionsPlayed         int
	TotalSessionsWon            int
	TotalTripleKills            int `json:",omitempty"`
	TotalTurretsKilled          int `json:",omitempty"`
	TotalUnrealKills            int `json:",omitempty"`
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
	Assists                         int   `json:",omitempty"`
	BarracksKilled                  int   `json:",omitempty"` // Number of enemy inhibitors killed.`json:",omitempty"`
	ChampionsKilled                 int   `json:",omitempty"`
	CombatPlayerScore               int   `json:",omitempty"`
	ConsumablesPurchased            int   `json:",omitempty"`
	DamageDealtPlayer               int   `json:",omitempty"`
	DoubleKills                     int   `json:",omitempty"`
	FirstBlood                      int   `json:",omitempty"`
	Gold                            int   `json:",omitempty"`
	GoldEarned                      int   `json:",omitempty"`
	GoldSpent                       int   `json:",omitempty"`
	Item0                           int64 `json:",omitempty"`
	Item1                           int64 `json:",omitempty"`
	Item2                           int64 `json:",omitempty"`
	Item3                           int64 `json:",omitempty"`
	Item4                           int64 `json:",omitempty"`
	Item5                           int64 `json:",omitempty"`
	Item6                           int64 `json:",omitempty"`
	ItemsPurchased                  int   `json:",omitempty"`
	KillingSprees                   int   `json:",omitempty"`
	LargestCriticalStrike           int   `json:",omitempty"`
	LargestKillingSpree             int   `json:",omitempty"`
	LargestMultiKill                int   `json:",omitempty"`
	LegendaryItemsCreated           int   `json:",omitempty"` // Number of tier 3 items built.`json:",omitempty"`
	Level                           int   `json:",omitempty"`
	MagicDamageDealtPlayer          int   `json:",omitempty"`
	MagicDamageDealtToChampions     int   `json:",omitempty"`
	MagicDamageTaken                int   `json:",omitempty"`
	MinionsDenied                   int   `json:",omitempty"`
	MinionsKilled                   int   `json:",omitempty"`
	NeutralMinionsKilled            int   `json:",omitempty"`
	NeutralMinionsKilledEnemyJungle int   `json:",omitempty"`
	NeutralMinionsKilledYourJungle  int   `json:",omitempty"`
	NexusKilled                     bool  `json:",omitempty"` // Flag specifying if the summoner got the killing blow on the nexus.`json:",omitempty"`
	NodeCapture                     int   `json:",omitempty"`
	NodeCaptureAssist               int   `json:",omitempty"`
	NodeNeutralize                  int   `json:",omitempty"`
	NodeNeutralizeAssist            int   `json:",omitempty"`
	NumDeaths                       int   `json:",omitempty"`
	NumItemsBought                  int   `json:",omitempty"`
	ObjectivePlayerScore            int   `json:",omitempty"`
	PentaKills                      int   `json:",omitempty"`
	PhysicalDamageDealtPlayer       int   `json:",omitempty"`
	PhysicalDamageDealtToChampions  int   `json:",omitempty"`
	PhysicalDamageTaken             int   `json:",omitempty"`
	QuadraKills                     int   `json:",omitempty"`
	SightWardsBought                int   `json:",omitempty"`
	Spell1Cast                      int   `json:",omitempty"` // Number of times first champion spell was cast.`json:",omitempty"`
	Spell2Cast                      int   `json:",omitempty"` // Number of times second champion spell was cast.`json:",omitempty"`
	Spell3Cast                      int   `json:",omitempty"` // Number of times third champion spell was cast.`json:",omitempty"`
	Spell4Cast                      int   `json:",omitempty"` // Number of times fourth champion spell was cast.`json:",omitempty"`
	SummonSpell1Cast                int   `json:",omitempty"`
	SummonSpell2Cast                int   `json:",omitempty"`
	SuperMonsterKilled              int   `json:",omitempty"`
	Team                            int   `json:",omitempty"`
	TeamObjective                   int   `json:",omitempty"`
	TimePlayed                      int   `json:",omitempty"`
	TotalDamageDealt                int   `json:",omitempty"`
	TotalDamageDealtToChampions     int   `json:",omitempty"`
	TotalDamageTaken                int   `json:",omitempty"`
	TotalHeal                       int   `json:",omitempty"`
	TotalPlayerScore                int   `json:",omitempty"`
	TotalScoreRank                  int   `json:",omitempty"`
	TotalTimeCrowdControlDealt      int   `json:",omitempty"`
	TotalUnitsHealed                int   `json:",omitempty"`
	TripleKills                     int   `json:",omitempty"`
	TrueDamageDealtPlayer           int   `json:",omitempty"`
	TrueDamageDealtToChampions      int   `json:",omitempty"`
	TrueDamageTaken                 int   `json:",omitempty"`
	TurretsKilled                   int   `json:",omitempty"`
	UnrealKills                     int   `json:",omitempty"`
	VictoryPointTotal               int   `json:",omitempty"`
	VisionWardsBought               int   `json:",omitempty"`
	WardKilled                      int   `json:",omitempty"`
	WardPlaced                      int   `json:",omitempty"`
	Win                             bool  //Flag specifying whether or not this game was won.
}

type Image struct {
	Full   string
	Group  string
	H      int
	Sprite string
	W      int
	X      int
	Y      int
}

func (i Image) GetImageURL() string {
	return fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/4.2.6/img/%s/%s", i.Group, i.Full)
}

type Info struct {
	Attack     int
	Defense    int
	Difficulty int
	Magic      int
}

type Passive struct {
	Description          string
	Image                Image
	Name                 string
	SanitizedDescription string
}

type Recommended struct {
	Blocks   []Block
	Champion string
	Map      string
	Mode     string
	Priority bool
	Title    string
	Type     string
}

type Stats struct {
	Armor                float64
	Armorperlevel        float64
	Attackdamage         float64
	Attackdamageperlevel float64
	Attackrange          float64
	Attackspeedoffset    float64
	Attackspeedperlevel  float64
	Crit                 float64
	Critperlevel         float64
	Hp                   float64
	Hpperlevel           float64
	Hpregen              float64
	Hpregenperlevel      float64
	Movespeed            float64
	Mp                   float64
	Mpperlevel           float64
	Mpregen              float64
	Mpregenperlevel      float64
	Spellblock           float64
	Spellblockperlevel   float64
}

type Block struct {
	Items   []BlockItem
	RecMath bool
	Type    string
}

type BlockItem struct {
	Count int
	Id    int64
}

type Skin struct {
	Id   int64
	Name string
	Num  int
}
type SpellVars struct {
	Coeff     []float64
	Dyn       string
	Key       string
	Link      string
	RanksWith string
}

type LevelTip struct {
	Effect []string
	Label  []string
}

type ItemList struct {
	Basic     BasicData
	Data      map[string]Item
	ItemsById map[int64]Item
	Groups    []Group
	Tree      []ItemTree
	Type      string
	Version   string
}

type Item struct {
	Colloq               string
	ConsumeOnFull        bool
	Consumed             bool
	Depth                int
	Description          string
	From                 []string
	Gold                 Gold
	Group                string
	HideFromAll          bool
	Id                   int64
	Image                Image
	InStore              bool
	Into                 []string
	Maps                 map[string]bool
	Name                 string
	Plaintext            string
	RequiredChampion     string
	Rune                 MetaData
	SanitizedDescription string
	SpecialRecipe        int
	Stacks               int
	Stats                BasicDataStats
	Tags                 []string
}

type ItemTree struct {
	Header string
	Tags   []string
}

type Group struct {
	MaxGroupOwnable string
	Key             string
}

type BasicData struct {
	Colloq               string
	ConsumeOnFull        bool
	Consumed             bool
	Depth                int
	Description          string
	From                 []string
	Gold                 Gold
	Group                string
	HideFromAll          bool
	Id                   int64
	Image                Image
	InStore              bool
	Into                 []string
	Maps                 map[string]bool
	Name                 string
	Plaintext            string
	RequiredChampion     string
	Rune                 MetaData
	SanitizedDescription string
	SpecialRecipe        int
	Stacks               int
	Stats                BasicDataStats
	Tags                 []string
}

type BasicDataStats struct {
	FlatArmorMod                        float64
	FlatAttackSpeedMod                  float64
	FlatBlockMod                        float64
	FlatCritChanceMod                   float64
	FlatCritDamageMod                   float64
	FlatEXPBonus                        float64
	FlatEnergyPoolMod                   float64
	FlatEnergyRegenMod                  float64
	FlatHPPoolMod                       float64
	FlatHPRegenMod                      float64
	FlatMPPoolMod                       float64
	FlatMPRegenMod                      float64
	FlatMagicDamageMod                  float64
	FlatMovementSpeedMod                float64
	FlatPhysicalDamageMod               float64
	FlatSpellBlockMod                   float64
	PercentArmorMod                     float64
	PercentAttackSpeedMod               float64
	PercentBlockMod                     float64
	PercentCritChanceMod                float64
	PercentCritDamageMod                float64
	PercentDodgeMod                     float64
	PercentEXPBonus                     float64
	PercentHPPoolMod                    float64
	PercentHPRegenMod                   float64
	PercentLifeStealMod                 float64
	PercentMPPoolMod                    float64
	PercentMPRegenMod                   float64
	PercentMagicDamageMod               float64
	PercentMovementSpeedMod             float64
	PercentPhysicalDamageMod            float64
	PercentSpellBlockMod                float64
	PercentSpellVampMod                 float64
	RFlatArmorModPerLevel               float64
	RFlatArmorPenetrationMod            float64
	RFlatArmorPenetrationModPerLevel    float64
	RFlatCritChanceModPerLevel          float64
	RFlatCritDamageModPerLevel          float64
	RFlatDodgeMod                       float64
	RFlatDodgeModPerLevel               float64
	RFlatEnergyModPerLevel              float64
	RFlatEnergyRegenModPerLevel         float64
	RFlatGoldPer10Mod                   float64
	RFlatHPModPerLevel                  float64
	RFlatHPRegenModPerLevel             float64
	RFlatMPModPerLevel                  float64
	RFlatMPRegenModPerLevel             float64
	RFlatMagicDamageModPerLevel         float64
	RFlatMagicPenetrationMod            float64
	RFlatMagicPenetrationModPerLevel    float64
	RFlatMovementSpeedModPerLevel       float64
	RFlatPhysicalDamageModPerLevel      float64
	RFlatSpellBlockModPerLevel          float64
	RFlatTimeDeadMod                    float64
	RFlatTimeDeadModPerLevel            float64
	RPercentArmorPenetrationMod         float64
	RPercentArmorPenetrationModPerLevel float64
	RPercentAttackSpeedModPerLevel      float64
	RPercentCooldownMod                 float64
	RPercentCooldownModPerLevel         float64
	RPercentMagicPenetrationMod         float64
	RPercentMagicPenetrationModPerLevel float64
	RPercentMovementSpeedModPerLevel    float64
	RPercentTimeDeadMod                 float64
	RPercentTimeDeadModPerLevel         float64
}

type Gold struct {
	Base        int
	Purchasable bool
	Sell        int
	Total       int
}

type MetaData struct {
	IsRune bool
	Tier   string
	Type   string
}

type SummonerSpell struct {
	Cooldown             []float64
	CooldownBurn         string
	Cost                 []int
	CostBurn             string
	CostType             string
	Description          string
	Effect               [][]int // This field is a List of List of Integer.
	EffectBurn           []string
	Id                   int64
	Image                Image
	Key                  string
	Leveltip             LevelTip
	Maxrank              int
	Modes                []string
	Name                 string
	Range                interface{} //object	This field is either a List of Integer or the String 'self' for spells that target one's own champion.
	RangeBurn            string
	Resource             string
	SanitizedDescription string
	SanitizedTooltip     string
	SummonerLevel        int
	Tooltip              string
	Vars                 []SpellVars
}
type ErrorStatus struct {
	Status StatusMessage
}

type StatusMessage struct {
	Message    string
	StatusCode int `json:"status_code"`
}
