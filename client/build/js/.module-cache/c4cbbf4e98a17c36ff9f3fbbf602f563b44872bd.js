define(['jquery'], function ($) {
	
	function getMatchHistory (name) {
		$.get('/api/summoner/matchHistory?name=' + name, function () {
			
		});
	}

	return {

	}
});