define(['jquery', './apiPaths.js'], function ($, apiPaths) {
	
	function getMatchHistory (name, callback) {
		$.get(apiPaths.SUMMONER.GET_MATCH_HISTORY(name), function (data) {
			var parseData = data;

			if(!parseData.error) {
				callback (parseData);	
			}
		});
	}

	function getRankedData (name, callback) {
		$.get(apiPaths.SUMMONER.GET_RANKED_DATA(name), function (data) {
			var parseData = data;

			if(!parseData.error) {
				callback (parseData);	
			}
		});
	}

	return {
		getMatchHistory: getMatchHistory,
		getRankedData: getRankedData
	}
	
});
