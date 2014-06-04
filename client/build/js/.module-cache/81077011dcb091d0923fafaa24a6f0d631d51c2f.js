define(['/js/settings.js', '/js/api/summonerAPI.js', 'js/mockServer/mockSummonerAPI.js', '/js/dispatcher.js'], function (settings, realApi, mockAPI, dispatcher) {
	var matchHistory = {};


	var api = settings.useLocal ? mockAPI : realApi;

	summonerStore.prototype = {
		addChangeListener: function (callback) {
			changeListeners.push(callback);
		},

		removeChangeListener: function (callback) {
			changeListeners = changeListeners.filter(function (a) {
				return a !== callback;
			});
		},

		notifyListeners: function () {
			changeListeners.forEach(function (a) {
				a();
			});
		},

		getMatchHistory: function (name) {
			if(matchHistory[name] === undefined) {
				matchHistory[name] = api.getMatchHistory(name);
			}
		}
	}

	return summonerStore;
})