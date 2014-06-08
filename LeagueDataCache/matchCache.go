package LeagueDataCache

import (
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

type champFetchFunc func(id int64, api *lapi.LolFetcher) (lapi.Champion, error)

func convertGamesToMatchHistory(id int64, games []lapi.Game) (MatchHistory, error) {
	summary := MatchHistory{SummonerId: id}
	for _, game := range games {
		lg := NewMatchSimpleFromGame(game)
		if champ, ok := allChampions[game.ChampionId]; ok {
			lg.ChampionName = champ.Name
			if champ.Id > 0 {
				lg.ChampionImage = champ.Image.GetImageURL()
			}
		}
		summary.Games = append(summary.Games, lg)
	}
	return summary, nil
}

// Fetches a cached match and returns detailed match.
func convertGameToMatchDetail(g lapi.Game) (MatchDetail, error) {
	lmd := NewMatchDetailsFromGame(g)
	if champ, ok := allChampions[g.ChampionId]; ok {
		lmd.ChampionName = champ.Name
	}
	for ind, p := range lmd.FellowPlayers {

		if champ, ok := allChampions[p.ChampionId]; ok {
			p.ChampionName = champ.Name
		}
		if summ, ok := allSummonersById[p.SummonerId]; ok {
			p.SummonerName = summ.Name
		}
		lmd.FellowPlayers[ind] = p
	}
	return lmd, nil
}
