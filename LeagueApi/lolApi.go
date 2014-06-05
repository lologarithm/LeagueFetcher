package LeagueApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TierChallenger = 10
	TierDiamond    = 5
	TierPlatinum   = 4
	TierGold       = 3
	TierSilver     = 2
	TierBronze     = 1
)

const (
	baseUrl         = "https://prod.api.pvp.net/api/lol"
	region          = "na"
	summonerVersion = "v1.4"
	statsVersion    = "v1.3"
	champVersion    = "v1.2"
	leagueVersion   = "v2.4"
	teamVersion     = "v2.3"
	gameVersion     = "v1.3"
)

type RemoteGet func(url string) (*http.Response, error)
type Logger interface {
	Infof(string, ...interface{})
}

var ApiKey string

type LolFetcher struct {
	Get RemoteGet
	Log Logger
}

func (lf *LolFetcher) makeUrl(version string, method string) string {
	url := fmt.Sprintf("%s/%s/%s/%s?api_key=%s", baseUrl, region, version, method, ApiKey)
	//lf.Log("URL: %s\n", url)
	return url
}

func (lf *LolFetcher) makeStaticDataUrl(version string, method string, params string) string {
	url := fmt.Sprintf("%s/static-data/%s/%s/%s?api_key=%s%s", baseUrl, region, version, method, ApiKey, params)
	//lf.Log("URL: %s\n", url)
	return url
}

func (lf *LolFetcher) makeRequest(url string, value interface{}) {
	st := time.Now().UnixNano()
	resp, err := lf.Get(url)
	if err != nil {
		lf.Log.Infof("Failed to open conn: %s\n", err.Error())
		return
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		lf.Log.Infof("Failed to open conn: %s\n", err.Error())
		return
	}

	unmarshErr := json.Unmarshal(body, value)
	if unmarshErr != nil {
		lf.Log.Infof("Failed to unmarshal json: %s\n", unmarshErr.Error())
	}
	lf.Log.Infof("Request (%s) Took: %.4fms\n", url, (float64(time.Now().UnixNano()-st))/float64(1000000.0))
}

func (lf *LolFetcher) GetSummonerByName(name string) (summoners map[string]Summoner) {
	name = NormalizeString(name)
	lf.makeRequest(lf.makeUrl(summonerVersion, "summoner/by-name/"+name), &summoners)
	return
}

func (lf *LolFetcher) GetSummonersById(ids []int64) (summoners map[string]Summoner) {
	var buffer bytes.Buffer
	buffer.WriteString("summoner/")
	for _, id := range ids {
		buffer.WriteString(strconv.FormatInt(id, 10))
		buffer.WriteString(",")
	}
	lf.makeRequest(lf.makeUrl(summonerVersion, buffer.String()), &summoners)
	return
}

func (lf *LolFetcher) GetSummonerRankedStats(id int64) (srs RankedStats) {
	method := fmt.Sprintf("stats/by-summoner/%d/ranked", id)
	lf.makeRequest(lf.makeUrl(statsVersion, method), &srs)
	return
}

func (lf *LolFetcher) GetSummonerSummaryStats(id int64) (stats PlayerStatsSummaryList) {
	method := fmt.Sprintf("stats/by-summoner/%d/summary", id)
	lf.makeRequest(lf.makeUrl(statsVersion, method), &stats)
	return
}

func (lf *LolFetcher) GetSummonerLeagues(id int64) (leagues map[string][]League) {
	method := fmt.Sprintf("league/by-summoner/%d/entry", id)
	lf.makeRequest(lf.makeUrl(leagueVersion, method), &leagues)
	return
}

func (lf *LolFetcher) GetSummonerTeams(id int64, get RemoteGet) (teams map[string][]Team) {
	method := fmt.Sprintf("team/by-summoner/%d", id)
	lf.makeRequest(lf.makeUrl(teamVersion, method), &teams)
	return
}

func (lf *LolFetcher) GetAllChampions() (champs ChampionList) {
	params := "&champData=all"
	lf.makeRequest(lf.makeStaticDataUrl(champVersion, "champion", params), &champs)
	return
}

func (lf *LolFetcher) GetChampion(id int64) (champ Champion) {
	if id <= 0 {
		champ.Name = "Total"
		return
	}
	params := "&champData=all"
	method := fmt.Sprintf("champion/%d", id)
	lf.makeRequest(lf.makeStaticDataUrl(champVersion, method, params), &champ)
	return
}

func (lf *LolFetcher) GetRecentMatches(id int64) (r RecentGames) {
	method := fmt.Sprintf("game/by-summoner/%d/recent", id)
	lf.makeRequest(lf.makeUrl(gameVersion, method), &r)
	return
}

// -1 means tier 1 is worse, 0 means equal, 1 means tier 1 is better
func CompareRanked(tier1 string, div1 string, tier2 string, div2 string) int {
	tierMap := map[string]int{"BRONZE": TierBronze, "SILVER": TierSilver, "GOLD": TierGold, "PLATINUM": TierPlatinum, "DIAMOND": TierDiamond, "CHALLENGER": TierChallenger}
	if tierMap[tier1] > tierMap[tier2] {
		return 1
	} else if tierMap[tier1] < tierMap[tier2] {
		return -1
	}

	if LeagueDivisionToNumber(div1) < LeagueDivisionToNumber(div2) {
		return 1
	} else if LeagueDivisionToNumber(div1) > LeagueDivisionToNumber(div2) {
		return -1
	}

	return 0
}

// Not enough different divisions to need a real roman numeral translator.
func LeagueDivisionToNumber(div string) int {
	switch div {
	case "V":
		return 5
	case "IV":
		return 4
	case "III":
		return 3
	case "II":
		return 2
	case "I":
		return 1
	}
	// Return some very large value
	return 100
}

func NormalizeString(s string) string {
	s = strings.ToLower(s)
	return strings.Replace(s, " ", "", -1)
}
