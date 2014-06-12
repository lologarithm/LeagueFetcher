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
		// TODO: Make sure this is only called from cache thread?
		if summ, ok := allSummonersById[p.SummonerId]; ok {
			p.SummonerName = summ.Name
		}
		lmd.FellowPlayers[ind] = p
	}
	lmd.Items = make([]ItemDetail, 7)
	if lmd.Stats.Item0 > 0 {
		item := allItems.ItemsById[lmd.Stats.Item0]
		lmd.Items[0] = ItemDetail{Name: item.Name, ImageUrl: item.Image.GetImageURL()}
	}
	if lmd.Stats.Item1 > 0 {
		item := allItems.ItemsById[lmd.Stats.Item1]
		lmd.Items[1] = ItemDetail{Name: item.Name, ImageUrl: item.Image.GetImageURL()}
	}
	if lmd.Stats.Item2 > 0 {
		item := allItems.ItemsById[lmd.Stats.Item2]
		lmd.Items[2] = ItemDetail{Name: item.Name, ImageUrl: item.Image.GetImageURL()}
	}
	if lmd.Stats.Item3 > 0 {
		item := allItems.ItemsById[lmd.Stats.Item3]
		lmd.Items[3] = ItemDetail{Name: item.Name, ImageUrl: item.Image.GetImageURL()}
	}
	if lmd.Stats.Item4 > 0 {
		item := allItems.ItemsById[lmd.Stats.Item4]
		lmd.Items[4] = ItemDetail{Name: item.Name, ImageUrl: item.Image.GetImageURL()}
	}
	if lmd.Stats.Item5 > 0 {
		item := allItems.ItemsById[lmd.Stats.Item5]
		lmd.Items[5] = ItemDetail{Name: item.Name, ImageUrl: item.Image.GetImageURL()}
	}
	if lmd.Stats.Item6 > 0 {
		item := allItems.ItemsById[lmd.Stats.Item6]
		lmd.Items[6] = ItemDetail{Name: item.Name, ImageUrl: item.Image.GetImageURL()}
	}
	return lmd, nil
}
