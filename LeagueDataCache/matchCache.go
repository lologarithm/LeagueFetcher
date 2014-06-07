package LeagueDataCache

import (
	"appengine"
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
		summary.Games = append(summary.Games, lg)
	}
	return summary, nil
}

// Fetches a cached match and returns detailed match.
func convertGameToMatchDetail(g lapi.Game, request Request, api *lapi.LolFetcher) (MatchDetail, error) {
	lmd := NewMatchDetailsFromGame(g)
	champ, fErr := fetchAndCacheChampion(g.ChampionId, api)
	if fErr != nil {
		return MatchDetail{}, fErr
	}
	lmd.ChampionName = champ.Name
	missingIds := []int64{}
	for ind, p := range lmd.FellowPlayers {

		champ, fErr := fetchAndCacheChampion(p.ChampionId, api)
		if fErr != nil {
			return MatchDetail{}, fErr
		}

		p.ChampionName = champ.Name
		if summ, ok := allSummonersById[p.SummonerId]; ok {
			p.SummonerName = summ.Name
		} else {
			missingIds = append(missingIds, p.SummonerId)
		}

		lmd.FellowPlayers[ind] = p
	}

	if len(missingIds) > 0 {
		request.Context.Infof("Trying to get from %d summoner names from datastore.", len(missingIds))
		missingS, err := request.Persist.GetSummoners(missingIds)

		if err == nil {
			for _, data := range missingS {
				for ind, p := range lmd.FellowPlayers {
					if data.Id == p.SummonerId {
						p.SummonerName = data.Name
						lmd.FellowPlayers[ind] = p
						cacheSummoner(data)
						break
					}
				}
			}
		} else if me, ok := err.(appengine.MultiError); ok {
			for i, merr := range me {
				if merr == nil {
					for ind, p := range lmd.FellowPlayers {
						if missingIds[i] == p.SummonerId {
							p.SummonerName = missingS[i].Name
							lmd.FellowPlayers[ind] = p
							cacheSummoner(missingS[i])
							break
						}
					}
				}
			}
		}
	}

	return lmd, nil
}
