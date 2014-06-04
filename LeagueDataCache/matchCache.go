package LeagueDataCache

import (
	lapi "github.com/lologarithm/LeagueFetcher/LeagueApi"
)

// Fetches all recent games for a summoner and caches them.
func fetchAndCacheGames(id int64) {
	matches := lapi.GetRecentMatches(id)
	for _, match := range matches.Games {
		allGames[match.GameId] = match
		gamesBySummoner[id] = append([]lapi.Game{match}, gamesBySummoner[id]...)
	}
}

// Gets all recent games for a summoner and converts them to simple format.
func fetchSimpleMatchHistory(id int64) MatchHistory {
	summary := MatchHistory{SummonerId: id}
	if _, ok := gamesBySummoner[id]; !ok {
		gamesBySummoner[id] = []lapi.Game{}
		fetchAndCacheGames(id)
	}

	games := gamesBySummoner[id]
	for _, game := range games {
		champ := fetchAndCacheChampion(game.ChampionId)
		lg := NewMatchSimpleFromGame(game)
		lg.ChampionName = champ.Name
		summary.Games = append(summary.Games, lg)
	}

	return summary
}

// Fetches a cached match and returns detailed match.
func fetchMatch(id int64) MatchDetail {
	g := allGames[id]
	lmd := NewMatchDetailsFromGame(g)
	champ := fetchAndCacheChampion(g.ChampionId)
	lmd.ChampionName = champ.Name
	summIds := make([]int64, len(lmd.FellowPlayers))
	for ind, p := range lmd.FellowPlayers {
		champ := fetchAndCacheChampion(p.ChampionId)
		summIds[ind] = p.SummonerId
		p.ChampionName = champ.Name
		lmd.FellowPlayers[ind] = p
	}
	summoners := fetchSummonersById(summIds)
	for index, summ := range summoners {
		lmd.FellowPlayers[index].SummonerName = summ.Name
	}
	return lmd
}
