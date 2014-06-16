/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', 'js/components/controls/rankedStatsElement.js', 'js/utility/tiersAndDivisions.js'], 

function ($, React, moment, summonerStore, RankedStatsElement, tiersAndDivions) {
	var GeneralView = React.createClass({
		getInitialState: function () {
			return {
				dataUpdate: false
			};
		},
		
		render: function () {
			console.log(summonerStore.getRankedData(this.props.name));

			var data = summonerStore.getRankedData(this.props.name);
			if (this.props.name !== '' && data) {
				var soloqUrl = '';
				var soloqDisplay = '';
				var soloqLP = '';
				if (data.Solo5sLeague.Entries) {
					soloqUrl = 'imgs/leagues/' + data.Solo5sLeague.Tier + '_' + data.Solo5sLeague.Entries[0].Division + '.png';
					soloqDisplay = data.Solo5sLeague.Tier + ' ' + data.Solo5sLeague.Entries[0].Division;
					soloqLP = data.Solo5sLeague.Entries[0].LeaguePoints + ' LP';
				} else {
					soloqUrl = 'imgs/leagues/unknown.png';
					soloqDisplay = 'Unranked';
				}



				//Compute top KDA champs
				var sections = [];
				sections.push(this.renderTopChamps(data));



				return this.transferPropsTo(
						<div>
							<div className="flexContainer padding-l">
								{this.renderTeam(data, "RANKED_TEAM_3x3", "Ranked 3's")}
								<div className="flex1" style={{'text-align':'center'}}>
									<div style={{'font-size':'18px'}}>
										Solo Queue
									</div>
									<div>
										<img src={soloqUrl} style={{'max-width':'100%'}} />
									</div>
									<div>
										{soloqDisplay}
									</div>
									<div>
										<small>{soloqLP}</small>
									</div>
								</div>
								{this.renderTeam(data, "RANKED_TEAM_5x5", "Ranked 5's")}
							</div>
							{sections}
						</div>);
			} else {
				return <div></div>;
			}
		},

		renderTeam: function (data, queueStr, displayStr) {
			if(data.RankedTeamLeagues) {
				var keys = Object.keys(data.RankedTeamLeagues);
				var topTeam = keys.map(function (a) {
					return data.RankedTeamLeagues[a];
				}).filter(function (a) {
					return a.Queue === queueStr;
				}).sort(function (a,b) {
					return a.Tier === b.Tier ? 
										(tiersAndDivions.DIVISIONS[a.Entries[0].Division] > tiersAndDivions.DIVISIONS[b.Entries[0].Division] ? -1 : 1) : 
										(tiersAndDivions.TIERS[a.Tier] > tiersAndDivions.TIERS[b.Tier] ? -1 : 1);
				})[0];

				var fivesUrl = '';
				var fivesDisplay = '';
				if (topTeam) {
					fivesUrl = 'imgs/leagues/' + topTeam.Tier + '_' + topTeam.Entries[0].Division + '.png';
					fivesDisplay = topTeam.Tier + ' ' + topTeam.Entries[0].Division;
				} else {
					fivesUrl = 'imgs/leagues/unknown.png';
					fivesDisplay = 'Unranked';
				}

				return <div className="flex1" style={{'text-align':'center'}}>
							<div style={{'font-size':'18px'}}>
								{displayStr}
							</div>
							<div>
								<img src={fivesUrl} style={{'max-width':'100%'}} />
							</div>
							<div>
								{fivesDisplay}
							</div>
						</div>;
			} else {
				return null;
			}
		},

		renderBestThreesTeam: function (data) {

		},

		renderTopChamps: function (data) {
			if(data.Champions) {
				var renderedElm = data.Champions.filter(function (a){ 
					return a.Stats.TotalSessionsPlayed >= 5 && a.ChampionName !== 'Total';
				})
				.sort(this.sortByKDA)
				.filter(function (a,i) {
					return i < 5;
				}).map(function (a) {
					return <RankedStatsElement data={a} />;
				});

				return <div className="padding-l" style={{'border-top':'1px #ccc solid'}}>
							<h3 className="margin-bottom-l">Best Champions</h3>
							{renderedElm}
						</div>;
			} else {
				return null;
			}
		},

		sortByKDA: function (a,b) {
			return this.computeKDA(a) > this.computeKDA(b) ? -1 : 1;
		},

		computeKDA: function (a) {
			return Math.round((a.Stats.TotalChampionKills + a.Stats.TotalAssists) / (a.Stats.TotalDeathsPerSession || 1) * 100) / 100;
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

	return GeneralView;
})