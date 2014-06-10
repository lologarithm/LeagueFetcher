/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js'], function ($, React, moment, summonerStore) {
	var MatchDetails = React.createClass({
		render: function () {
			fellowPlayers = this.props.data.FellowPlayers || [];


			var bluePlayers = fellowPlayers.filter(function (a) {
				return a.Side === 'blue';
			});
			var purplePlayers = fellowPlayers.filter(function (a) {
				return a.Side === 'purple';
			});

			var html = bluePlayers.map($.proxy(function (a, index) {
				return 	<tr>
							<td className="padding-right-xl padding-left-m" style={{'background-color':'#8cd6f4', 'font-weight':(this.props.name === bluePlayers[index].SummonerName ? 'bold':'normal')}}>
								<a href="#" onClick={this.onClick(bluePlayers[index].SummonerName)}>{bluePlayers[index].SummonerName}</a>
							</td>
							<td className="padding-right-m" style={{'background-color':'#8cd6f4'}}>{bluePlayers[index].ChampionName}</td>
							<td className="padding-left-l padding-right-l"><span style={{'border-right':'1px #ccc solid'}}> </span></td>
							<td className="padding-right-xl padding-left-m" style={{'background-color':'#caa4cc', 'font-weight':(this.props.name === purplePlayers[index].SummonerName ? 'bold':'normal')}}>
								<a href="#" onClick={this.onClick(purplePlayers[index].SummonerName)}>{purplePlayers[index].SummonerName}</a>
							</td>
							<td className="padding-right-m" style={{'background-color':'#caa4cc'}}>{purplePlayers[index].ChampionName}</td>
						</tr>;
			}, this));

			return this.transferPropsTo(
						<div className='padding-m'>
							<table>
								<thead>
									<td><span>Gold Earned:{this.props.data.Stats.GoldEarned || 0}, Minions Killed:{this.props.data.Stats.MinionsKilled || 0}</span></td>
								</thead>
								{html}
							</table>
						</div>);
		},

		onClick: function (name) {
			return (function () {
				this.props.searchName(name);
			}).bind(this);
		}
	});

	return MatchDetails;
})