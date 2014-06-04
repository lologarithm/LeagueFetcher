/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', './matchHistoryElement.js'], function ($, React, moment, summonerStore, MatchHistoryElement) {
	var MatchHistoryPanel = React.createClass({
		render: function () {

			var items = summonerStore.getMatchHistory(this.props.name).Games.map(function (a) {
				return <MatchHistoryElement data={a} />;
			});

			return this.transferPropsTo(<div>{items}</div>);
		},

		componentDidMount: function ( ){
			summonerStore.addChangeListener(this.onChange);
		},

		componentWillUnmount: function () {
			summonerStore.removeChangeListener(this.onChange);
		},

		onChange: function () {
			this.setState({});
		}
	});

	return MatchHistoryPanel;
})