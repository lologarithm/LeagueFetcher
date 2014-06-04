define(['./mockMatchHistoryResponse.js'], function (matchHistory) {
	

	function getMatchHistory(name, callback) {
		return callback(matchHistory);
	}

	return { 
		getMatchHistory: getMatchHistory
	};
})