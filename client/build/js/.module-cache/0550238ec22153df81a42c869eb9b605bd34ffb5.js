define(['jquery'], function ($) {
	
	function getMatchHistory (name) {
		var waiting = true;
		var result;

		$.get('/api/summoner/matchHistory?name=' + name, function (data) {
			result = data;
		});

		while(waiting) {

		}

		return result;
	}

	return {
		getMatchHistory
	}
});