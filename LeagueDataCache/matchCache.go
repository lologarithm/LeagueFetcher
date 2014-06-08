package LeagueDataCache

import (
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

type champFetchFunc func(id int64, api *lapi.LolFetcher) (lapi.Champion, error)

func convertGamesToMatchHistory(id int64, games []lapi.Game, getChamp champFetchFunc, api *lapi.LolFetcher) (MatchHistory, error) {
	summary := MatchHistory{SummonerId: id}
	for _, game := range games {
		champ, fErr := getChamp(game.ChampionId, api)
		if fErr != nil {
			return summary, fErr
		}
		lg := NewMatchSimpleFromGame(game)
		lg.ChampionName = champ.Name
		if champ.Id > 0 {
			lg.ChampionImage = champ.Image.GetImageURL()
		}
		summary.Games = append(summary.Games, lg)
	}
	return summary, nil
}

// Fetches a cached match and returns detailed match.
func convertGameToMatchDetail(g lapi.Game, api *lapi.LolFetcher) (MatchDetail, error) {
	lmd := NewMatchDetailsFromGame(g)
	champ, fErr := fetchAndCacheChampion(g.ChampionId, api)
	if fErr != nil {
		return MatchDetail{}, fErr
	}
	lmd.ChampionName = champ.Name
	for ind, p := range lmd.FellowPlayers {

		champ, fErr := fetchAndCacheChampion(p.ChampionId, api)
		if fErr != nil {
			return MatchDetail{}, fErr
		}

		p.ChampionName = champ.Name
		if summ, ok := allSummonersById[p.SummonerId]; ok {
			p.SummonerName = summ.Name
		}
		lmd.FellowPlayers[ind] = p
	}
	return lmd, nil
}
