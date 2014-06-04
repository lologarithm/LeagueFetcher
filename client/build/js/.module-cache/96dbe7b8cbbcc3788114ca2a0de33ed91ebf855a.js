define(['/js/settings.js'], function (settings) {
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