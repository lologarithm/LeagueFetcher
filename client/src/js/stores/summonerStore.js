define(['/js/settings.js', '/js/api/summonerAPI.js', '/js/mockServer/mockSummonerAPI.js', '/js/dispatcher.js'], function (settings, realAPI, mockAPI, dispatcher) {
	var matchHistory = {};
	var rankedData = {};
	var changeListeners = [];

	var api = settings.useLocal ? mockAPI : realAPI;

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
				api.getMatchHistory(name, $.proxy(function (data) {
					matchHistory[name] = data;
					this.notifyListeners();
				}, this))
			}

			return matchHistory[name] || { SummonerId: 0, Games: [] };
		},

		getRankedData: function (name) {
			if(rankedData[name] === undefined) {
				api.getRankedData(name, $.proxy(function (data) {
					rankedData[name] = data;
					this.notifyListeners();
				}, this))
			}


			return rankedData[name] ||  { Id: 0, ProfileIconId: 0, SummonerLevel: 0, RevisionDate: 0, SummonerId: 0, Champions: [],
											ModifyDate: 0,
											Solo5sLeague: {
												Entries: [],
												Name: "",
												ParticipantId: "",
												Queue: "",
												Tier: ""
											},
											Solo3sLeague: {
												Entries: null,
												Name: "",
												ParticipantId: "",
												Queue: "",
												Tier: ""
											},
											RankedTeamLeagues: null,
											ExpireTime: 0
										};
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