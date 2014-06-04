/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', './matchHistoryPanel.js'], function ($, React, moment, summonerStore, MatchHistoryPanel) {
	var MatchHistoryPanel = React.createClass({displayName: 'MatchHistoryPanel',
		render: function () {

			var items = summonerStore.getMatchHistory(this.props.name).games.map(function (a) {
				return MatchHistoryElement( {data:a} );
			});

			return React.DOM.div(null, items);
		}
	});

	return MatchHistoryPanel;
})