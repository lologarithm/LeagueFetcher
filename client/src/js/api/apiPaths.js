define([], function () {
	return {
		SUMMONER : {
			GET_MATCH_HISTORY : function(name){return '/api/summoner/matchHistory?=' + name},
			GET_RANKED_DATA : function(name){return 'api/summoner/rankedData?name=' + name}
		},
		CHAMPION : {

		},
		MATCH : {
			GET_MATCH : function(id){ return '/api/match?id=' + id }
		}
	};
})