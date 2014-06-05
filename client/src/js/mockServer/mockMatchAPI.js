define(['./mockMatchResponse.js'], function (match) {
	

	function getMatch(id, callback) {
		return callback (match);
	}

	return { 
		getMatch: getMatch
	};
})