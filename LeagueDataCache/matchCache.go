package LeagueDataCache

import (
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

type champFetchFunc func(id int64, api *lapi.LolFetcher) lapi.Champion

func convertGamesToMatchHistory(id int64, games []lapi.Game, getChamp champFetchFunc, api *lapi.LolFetcher) MatchHistory {
	summary := MatchHistory{SummonerId: id}
	for _, game := range games {
		champ := getChamp(game.ChampionId, api)
		lg := NewMatchSimpleFromGame(game)
		lg.ChampionName = champ.Name
		summary.Games = append(summary.Games, lg)
	}
	return summary
}

// Fetches a cached match and returns detailed match.
func convertGameToMatchDetail(g lapi.Game, api *lapi.LolFetcher) MatchDetail {
	lmd := NewMatchDetailsFromGame(g)
	champ := fetchAndCacheChampion(g.ChampionId, api)
	lmd.ChampionName = champ.Name
	for ind, p := range lmd.FellowPlayers {
		champ := fetchAndCacheChampion(p.ChampionId, api)
		p.ChampionName = champ.Name
		if summ, ok := allSummonersById[p.SummonerId]; ok {
			p.SummonerName = summ.Name
		}
		lmd.FellowPlayers[ind] = p
	}
	return lmd
}
