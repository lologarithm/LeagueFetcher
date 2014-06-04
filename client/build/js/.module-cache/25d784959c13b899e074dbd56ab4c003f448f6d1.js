define(['/js/settings.js', '/js/api/summonerAPI.js', 'js/mockServer/mockSummonerAPI.js', '/js/dispatcher.js'], function (settings, realApi, mockAPI, dispatcher) {
	var matchHistory = {};


	var api = settings.useLocal ? mockAPI : realApi;

	var summonerStore = function () {

	};

	var summonerStore = {
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

			return matchHistory[name];
		}
	}

	//Register this store to react appropriately to UI events
	dispatcher.register(function (action) {
		switch(action.type) {
			case ModelActionTypes.LINK_ADDED:
				modelStore.addLink(action.start, action.end);
				break;
			case ModelActionTypes.LINK_REMOVED:
				modelStore.removeLink(action.id);
				break;
		}
	});

	return summonerStore;
})