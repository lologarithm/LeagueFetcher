define(['jquery', './apiPaths.js'], function ($, apiPaths) {
	
	function getMatchHistory (name, callback) {
		$.get(apiPaths.SUMMONER.GET_MATCH_HISTORY(name), function (data) {
			callback (JSON.parse(data));
		});
	}

	return {
		getMatchHistory: getMatchHistory
	}
	
});
