/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js'], function ($, React, moment, summonerStore) {
	var RankedStatsElement = React.createClass({
		render: function () {
			var elmStyle = {};
			var winRate = 'Win Rate: ' + Math.round(this.props.data.Stats.TotalSessionsWon / this.props.data.Stats.TotalSessionsPlayed * 100) + '%';
			var kda = 'KDA: ' + this.formatNumber((this.props.data.Stats.TotalChampionKills + this.props.data.Stats.TotalAssists) / this.props.data.Stats.TotalDeathsPerSession);

			return <div className='matchHistoryElement padding-m' style={elmStyle}  >
						<img style={{'width': '50px'}} src={this.props.data.ChampionImage} />
						<span style={{'paddingLeft': '15px', 'width':'100px','display':'inline-block'}}>{kda}</span>
						<span>
							<span style={{'color':'green'}}>{this.formatNumber(this.props.data.Stats.TotalChampionKills / this.props.data.Stats.TotalSessionsPlayed)}</span>/
							<span style={{'color':'red'}}>{this.formatNumber(this.props.data.Stats.TotalDeathsPerSession / this.props.data.Stats.TotalSessionsPlayed)}</span>/
							<span style={{'color':'orange'}}>{this.formatNumber(this.props.data.Stats.TotalAssists / this.props.data.Stats.TotalSessionsPlayed)}</span>
						</span>
						<span style={{'display':'inline-block','float':'right'}}>{winRate}</span>
					</div>;
		},

		formatNumber: function (num) {
			return Math.round(num / (this.props.data.Stats.NumDeaths || 1) * 100) / 100
		}
	});

	return RankedStatsElement;
})