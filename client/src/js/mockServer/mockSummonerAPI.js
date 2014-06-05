define(['./mockMatchHistoryResponse.js', './mockRankedDataResponse.js'], function (matchHistory, rankedData) {
	

	function getMatchHistory(name, callback) {
		return setTimeout(function () {
			callback (matchHistory);
		}, 250);
	}

	function getRankedData(name, callback) {
		return setTimeout(function () {
			callback (rankedData);
		}, 250);
	}

	return { 
		getMatchHistory: getMatchHistory,
		getRankedData: getRankedData
	};
})