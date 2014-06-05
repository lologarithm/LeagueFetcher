/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', 'js/components/controls/rankedStatsElement.js'], 

function ($, React, moment, summonerStore, RankedStatsElement) {
	var RankedStatsView = React.createClass({
		getInitialState: function () {
			return {
				dataUpdate: false
			};
		},
		
		render: function () {
			if (this.props.name !== '') {
				if(summonerStore.getRankedData(this.props.name).Champions && summonerStore.getRankedData(this.props.name).Champions.length > 0) {
					var elements = summonerStore.getRankedData(this.props.name).Champions.sort(this.sortByMostPlayed).map(function (a) {
						return <RankedStatsElement data={a} />;
					})


					return <div style={{'width':'500px', 'margin':'10px auto 0'}} >
								<h3>Ranked Stats</h3>
								<div className="margin-top-m">
									{elements}
								</div>
							</div>;
				} else {
					return <div style={{'width':'500px', 'margin':'10px auto 0'}}>
								<h3>Ranked Stats</h3>
								<span>No ranked games found</span>
							</div>;
				}
			} else {
				return <div></div>;
			}
		},

		sortByMostPlayed: function (a,b) {
			return a.Stats.TotalSessionsPlayed > b.Stats.TotalSessionsPlayed ? -1 : 1;
		},

		componentDidMount: function ( ){
			summonerStore.addChangeListener(this.onChange);
		},

		componentWillUnmount: function () {
			summonerStore.removeChangeListener(this.onChange);
		},

		componentDidUpdate: function () {
			this.state.dataUpdate = false;
		},

		shouldComponentUpdate: function (nextProps) {
			return this.props.name !== nextProps.name || this.state.dataUpdate;
		},

		onChange: function () {
			this.setState({ dataUpdate : true });
		}
	});

	return RankedStatsView;
})