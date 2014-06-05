/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', './matchHistoryElement.js', './popover.js'], 

function ($, React, moment, summonerStore, MatchHistoryElement, Popover) {
	var MatchHistoryPanel = React.createClass({
		getInitialState: function () {
			return {
				target: null
			}
		},

		render: function () {
			if(summonerStore.getMatchHistory(this.props.name).Games && summonerStore.getMatchHistory(this.props.name).Games.length > 0) {
				var popover = summonerStore.getMatchHistory(this.props.name).Games.filter($.proxy(function (a) {
					this.state.target === a.GameId;
				}, this)).map(function (a) {
					return <Popover>a.GameId</Popover>;
				})[0];

				var items = summonerStore.getMatchHistory(this.props.name).Games.map(function (a) {
					return <MatchHistoryElement data={a} />;
				});

				return this.transferPropsTo(<div>{items}</div>);
			} else {
				return this.transferPropsTo(<div>No recent games found</div>);
			}
			
		},

		componentDidMount: function ( ){
			summonerStore.addChangeListener(this.onChange);
		},

		componentWillUnmount: function () {
			summonerStore.removeChangeListener(this.onChange);
		},

		onChange: function () {
			this.setState({});
		},

		elementClickd: function () {
			this.setState({ target : a.GameId });
		}
	});

	return MatchHistoryPanel;
})