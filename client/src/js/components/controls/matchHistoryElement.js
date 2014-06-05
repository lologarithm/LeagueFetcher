/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js'], function ($, React, moment, summonerStore) {
	var MatchHistoryPanel = React.createClass({
		render: function () {
			var kda = Math.round((this.props.data.Stats.ChampionsKilled + this.props.data.Stats.Assists) / (this.props.data.Stats.NumDeaths || 1) * 100) / 100;
			var elmStyle = this.props.data.Stats.Win ? {'border-top':'2px solid green'} : {'border-top':'2px solid red'};

			return <div className='matchHistoryElement padding-m margin-bottom-m' style={elmStyle} >
						<div>
							<span>Champion Played: {this.props.data.ChampionName}</span>
							<span style={{'display':'inline-block', 'float':'right'}}>
								<span style={{'color':'green'}}>{this.props.data.Stats.ChampionsKilled}</span>/
								<span style={{'color':'red'}}>{this.props.data.Stats.NumDeaths}</span>/
								<span style={{'color':'orange'}}>{this.props.data.Stats.Assists}</span>
							</span>
						</div>
						<div>
							<span>Date Played: { moment(this.props.data.CreateDate).format('MMMM Do YYYY, h:mm a') }</span>
							<span style={{'display':'inline-block', 'float':'right'}}>KDA: {kda}</span>
						</div>
						<div>Game Mode: { this.props.data.GameMode }</div>
						<div>IP Earned: { this.props.data.IpEarned }</div>
					</div>;
		}
	});

	return MatchHistoryPanel;
})