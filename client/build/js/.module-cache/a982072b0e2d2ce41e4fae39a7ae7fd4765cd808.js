define(['/js/settings.js', '/js/api/summonerAPI.js', 'js/mockServer/mockSummonerAPI.js'], function (settings) {
	var matchHistory = {};

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

		getMatchHistory: function () {

		}
	}

	return summonerStore;
})