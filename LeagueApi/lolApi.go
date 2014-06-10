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
	staticVersion   = "v1.2"
	leagueVersion   = "v2.4"
	teamVersion     = "v2.3"
	gameVersion     = "v1.3"
)

type RemoteGet func(url string) (*http.Response, error)
type Logger interface {
	Infof(string, ...interface{})
	Warningf(string, ...interface{})
}

var ApiKey string

type LolFetcher struct {
	Get RemoteGet
	Log Logger
}

type ApiAsyncResponse struct {
	Value interface{}
	Error *FetchError
}

type FetchError struct {
	Message string
	Code    int
}

func (fe *FetchError) Error() string {
	return fmt.Sprintf("Error: %s, Code: %d", fe.Message, fe.Code)
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

func (lf *LolFetcher) makeRequest(url string, value interface{}) *FetchError {
	st := time.Now().UnixNano()
	resp, err := lf.Get(url)
	if err != nil {
		lf.Log.Warningf("Failed to open conn: %s\n", err.Error())
		return &FetchError{Message: "Connection Failed.", Code: 0}
	}

	if resp.StatusCode != 200 {
		lf.Log.Warningf("Request (%s) failed with code %d, Status: %s", url, resp.StatusCode, resp.Status)
		return &FetchError{Message: resp.Status, Code: resp.StatusCode}
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		lf.Log.Warningf("Failed to parse response body.: %s\n", err.Error())
		return &FetchError{Message: err.Error(), Code: resp.StatusCode}
	}

	unmarshErr := json.Unmarshal(body, value)
	if unmarshErr != nil {
		lf.Log.Warningf("Failed to unmarshal json: %s\n", unmarshErr.Error())
		lf.Log.Infof("Failed JSON: %s", string(body))
		var sErr ErrorStatus
		statsError := json.Unmarshal(body, &sErr)
		if statsError == nil {
			lf.Log.Warningf("Rate Limit Exceeded on request: %s", url)
			return &FetchError{Message: sErr.Status.Message, Code: sErr.Status.StatusCode}
		}
	}

	lf.Log.Infof("LOLAPI Fetch (%s) Took: %.4fms\n", url, (float64(time.Now().UnixNano()-st))/float64(1000000.0))
	return nil
}

func (lf *LolFetcher) GetSummonerByName(name string) (summoners map[string]Summoner, limit *FetchError) {
	name = NormalizeString(name)
	limit = lf.makeRequest(lf.makeUrl(summonerVersion, "summoner/by-name/"+name), &summoners)
	return
}

func (lf *LolFetcher) GetSummonersById(ids []int64) (summoners map[string]Summoner, limit *FetchError) {
	var buffer bytes.Buffer
	buffer.WriteString("summoner/")
	for _, id := range ids {
		buffer.WriteString(strconv.FormatInt(id, 10))
		buffer.WriteString(",")
	}
	limit = lf.makeRequest(lf.makeUrl(summonerVersion, buffer.String()), &summoners)
	return
}

func (lf *LolFetcher) GetSummonerRankedStats(id int64) (srs RankedStats, limit *FetchError) {
	method := fmt.Sprintf("stats/by-summoner/%d/ranked", id)
	limit = lf.makeRequest(lf.makeUrl(statsVersion, method), &srs)
	return
}

func (lf *LolFetcher) GetSummonerRankedStatsAsync(id int64, resp chan ApiAsyncResponse) {
	var srs RankedStats
	method := fmt.Sprintf("stats/by-summoner/%d/ranked", id)
	limit := lf.makeRequest(lf.makeUrl(statsVersion, method), &srs)
	val := ApiAsyncResponse{Value: srs, Error: limit}
	resp <- val
}

func (lf *LolFetcher) GetSummonerSummaryStats(id int64) (stats PlayerStatsSummaryList, e *FetchError) {
	method := fmt.Sprintf("stats/by-summoner/%d/summary", id)
	e = lf.makeRequest(lf.makeUrl(statsVersion, method), &stats)
	return
}

func (lf *LolFetcher) GetSummonerLeagues(id int64) (leagues map[string][]League, e *FetchError) {
	method := fmt.Sprintf("league/by-summoner/%d/entry", id)
	e = lf.makeRequest(lf.makeUrl(leagueVersion, method), &leagues)
	return
}

func (lf *LolFetcher) GetSummonerLeaguesAsync(id int64, resp chan ApiAsyncResponse) {
	var leagues map[string][]League
	method := fmt.Sprintf("league/by-summoner/%d/entry", id)
	e := lf.makeRequest(lf.makeUrl(leagueVersion, method), &leagues)
	val := ApiAsyncResponse{Value: leagues, Error: e}
	resp <- val
}

func (lf *LolFetcher) GetSummonerTeams(id int64, get RemoteGet) (teams map[string][]Team, e *FetchError) {
	method := fmt.Sprintf("team/by-summoner/%d", id)
	e = lf.makeRequest(lf.makeUrl(teamVersion, method), &teams)
	return
}

func (lf *LolFetcher) GetAllChampions() (champs ChampionList, e *FetchError) {
	params := "&champData=all"
	e = lf.makeRequest(lf.makeStaticDataUrl(staticVersion, "champion", params), &champs)
	return
}

func (lf *LolFetcher) GetAllItems() (items ItemList, e *FetchError) {
	params := "&itemListData=all"
	e = lf.makeRequest(lf.makeStaticDataUrl(staticVersion, "item", params), &items)
	return
}

func (lf *LolFetcher) GetChampion(id int64) (champ Champion, e *FetchError) {
	if id <= 0 {
		champ.Name = "Total"
		return
	}
	params := "&champData=all"
	method := fmt.Sprintf("champion/%d", id)
	e = lf.makeRequest(lf.makeStaticDataUrl(staticVersion, method, params), &champ)
	return
}

func (lf *LolFetcher) GetRecentMatches(id int64) (r RecentGames, e *FetchError) {
	method := fmt.Sprintf("game/by-summoner/%d/recent", id)
	e = lf.makeRequest(lf.makeUrl(gameVersion, method), &r)
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
