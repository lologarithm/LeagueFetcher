define(['./mockMatchResponse.js'], function (match) {
	

	function getMatch(id, callback) {
		return setTimeout(function () {
			callback (match);
		}, 250);
	}

	return { 
		getMatch: getMatch
	};
})