define([], function () {
	return {
		SUMMONER : {
			GET_MATCH_HISTORY : function(name){return '/api/summoner/matchHistory?name=' + name},
			GET_RANKED_DATA : function(name){return 'api/summoner/rankedData?name=' + name}
		},
		CHAMPION : {

		},
		MATCH : {
			GET_MATCH : function(matchId, summonerId){ return '/api/match?matchId=' + matchId + '&summonerId=' + summonerId }
		}
	};
})