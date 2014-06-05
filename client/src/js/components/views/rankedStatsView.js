/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', 'js/components/controls/rankedStatsElements.js'], 

function ($, React, moment, summonerStore, RankedStatsElement) {
	var RankedStatsView = React.createClass({
		
		render: function () {
			if (this.props.name !== '') {
				var elements = summonerStore.getRankedData().Champions.map(function (a) {
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
	});

	return RankedStatsView;
})