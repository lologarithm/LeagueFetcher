define(['./mockMatchHistoryResponse.js'], function (matchHistory) {
	

	function getMatchHistory(name) {
		return matchHistory;
	}

	return { 
		getMatchHistory: getMatchHistory
	};
})