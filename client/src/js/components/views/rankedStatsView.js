/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', 'js/components/controls/rankedStatsElement.js'], 

function ($, React, moment, summonerStore, RankedStatsElement) {
	var RankedStatsView = React.createClass({
		
		render: function () {
			if (this.props.name !== '') {
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

		onChange: function () {
			this.setState({});
		}
	});

	return RankedStatsView;
})