/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', './matchHistoryElement.js'], function ($, React, moment, summonerStore, MatchHistoryElement) {
	var MatchHistoryPanel = React.createClass({displayName: 'MatchHistoryPanel',
		render: function () {

			var items = summonerStore.getMatchHistory(this.props.name).Games.map(function (a) {
				return MatchHistoryElement( {data:a} );
			});

			return this.transferPropsTo(React.DOM.div(null, items));
		}
	});

	return MatchHistoryPanel;
})