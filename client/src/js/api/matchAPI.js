define(['jquery', './apiPaths.js'], function ($, apiPaths) {
	
	function getMatch (matchId, summonerId, callback) {
		$.get(apiPaths.MATCH.GET_MATCH(matchId, summonerId), function (data) {
			callback (JSON.parse(data));
		});
	}

	return {
		getMatch: getMatch
	}
});