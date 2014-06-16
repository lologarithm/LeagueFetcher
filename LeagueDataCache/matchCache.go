package LeagueDataCache

import (
	"errors"

	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

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

func convertCachedMatchDetailsToHistory(games []cachedMatchDetail) (MatchHistory, error) {
	if len(games) == 0 {
		return MatchHistory{}, errors.New("No games to convert.")
	}
	summary := MatchHistory{SummonerId: games[0].SummonerId}
	for _, game := range games {
		ms, err := game.ToMatchSimple()
		if err != nil {
			return summary, err
		}
		summary.Games = append(summary.Games, ms)
	}
	return summary, nil
}

func convertMatchDetailsToHistory(id int64, games []MatchDetail) (MatchHistory, error) {
	summary := MatchHistory{SummonerId: id, Games: make([]MatchSimple, len(games))}
	if len(games) == 0 {
		return summary, errors.New("No games to convert.")
	}
	for ind, game := range games {
		summary.Games[ind] = game.toMatchSimple()
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
			if champ.Id > 0 {
				p.ChampionImage = champ.Image.GetImageURL()
			}
		} else if p.ChampionId == 0 {
			p.ChampionName = "Total"
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

func NewMatchSimpleFromGame(g lapi.Game) (lm MatchSimple) {
	lm.CreateDate = g.CreateDate
	lm.GameId = g.GameId
	lm.GameMode = g.GameMode
	lm.GameType = g.GameType
	lm.Invalid = g.Invalid
	lm.IpEarned = g.IpEarned
	lm.MapId = g.MapId
	if g.TeamId == 100 {
		lm.Side = "blue"
	} else {
		lm.Side = "purple"
	}
	lm.Stats = MatchStats{Assists: g.Stats.Assists, ChampionsKilled: g.Stats.ChampionsKilled, NumDeaths: g.Stats.NumDeaths, Win: g.Stats.Win}
	lm.SubType = g.SubType
	return
}

func NewMatchDetailsFromGame(g lapi.Game) (lmd MatchDetail) {
	lmd.CreateDate = g.CreateDate
	lmd.GameId = g.GameId
	lmd.GameMode = g.GameMode
	lmd.GameType = g.GameType
	lmd.Invalid = g.Invalid
	lmd.IpEarned = g.IpEarned
	lmd.MapId = g.MapId
	if g.TeamId == 100 {
		lmd.Side = "blue"
	} else {
		lmd.Side = "purple"
	}
	lmd.SubType = g.SubType
	players := []Player{}
	for _, player := range g.FellowPlayers {
		p := Player{ChampionId: player.ChampionId, SummonerId: player.SummonerId}
		if player.TeamId == 100 {
			p.Side = "blue"
		} else {
			p.Side = "purple"
		}
		players = append(players, p)
	}
	lmd.FellowPlayers = players
	lmd.Spell1 = g.Spell1
	lmd.Spell2 = g.Spell2
	lmd.Stats = g.Stats
	return
}
