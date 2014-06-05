define(['./mockMatchHistoryResponse.js', './mockRankedDataResponse.js'], function (matchHistory, rankedData) {
	

	function getMatchHistory(name, callback) {
		return callback (matchHistory);
	}

	function getRankedData(name, callback) {
		return callback (rankedData);
	}

	return { 
		getMatchHistory: getMatchHistory,
		getRankedData: getRankedData
	};
})