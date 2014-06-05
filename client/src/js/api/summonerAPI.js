define(['jquery', './apiPaths.js'], function ($, apiPaths) {
	
	function getMatchHistory (name, callback) {
		$.get(apiPaths.SUMMONER.GET_MATCH_HISTORY(name), function (data) {
			callback (JSON.parse(data));
		});
	}

	function getRankedData (name, callback) {
		$.get(apiPaths.SUMMONER.GET_RANKED_DATA(name), function (data) {
			callback (JSON.parse(data));
		});
	}

	return {
		getMatchHistory: getMatchHistory,
		getRankedData: getRankedData
	}
});