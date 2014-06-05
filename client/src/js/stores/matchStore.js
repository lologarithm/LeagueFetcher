define(['/js/settings.js', '/js/api/matchAPI.js', '/js/mockServer/mockMatchAPI.js', '/js/dispatcher.js'], function (settings, realAPI, mockAPI, dispatcher) {
	var matches = {};
	var changeListeners = [];

	var api = settings.useLocal ? mockAPI : realAPI;

	var matchStore = {
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

		getMatch: function (id, summonerId) {
			if(matches[id] === undefined) {
				api.getMatch(id, summonerId, $.proxy(function (data) {
					matches[id] = data;
					this.notifyListeners();
				}, this))
			}

			return matches[id] || { GameId: 0, FellowPlayers: [] };
		},
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

	return matchStore;
})