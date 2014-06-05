define(['jquery', './apiPaths.js'], function ($, apiPaths) {
	
	function getMatch (id, callback) {
		$.get(apiPaths.MATCH.GET_MATCH(id), function (data) {
			callback (JSON.parse(data));
		});
	}

	return {
		getMatch: getMatch
	}
});