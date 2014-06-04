define(['jquery'], function ($) {
	
	function getMatchHistory (name) {
		var waiting = true;
		var result;
		$.ajax({
	        url: '/api/summoner/matchHistory?name=' + name,
	        success: function (data) {
				result = JSON.parse(data);
			},
	        async:  false
	    });

		return result;
	}

	return {
		getMatchHistory: getMatchHistory
	}
	
});
