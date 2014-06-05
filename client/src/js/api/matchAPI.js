define(['jquery', './apiPaths.js'], function ($, apiPaths) {
	
	function getMatch (matchId, summonerId, callback) {
		$.get(apiPaths.MATCH.GET_MATCH(matchId, summonerId), function (data) {
			var parseData = JSON.parse(data);

			if(!parseData.error) {
				callback (parseData);	
			}
		});
	}

	return {
		getMatch: getMatch
	}
});