/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js'], function ($, React, moment, summonerStore) {
	var MatchDetails = React.createClass({
		render: function () {
			fellowPlayers = this.props.data.FellowPlayers || [];

			if(fellowPlayers.length > 0) {
				fellowPlayers = fellowPlayers.concat({ ChampionName: this.props.data.ChampionName, SummonerName: <span style={{'font-weight':'bold'}}>{this.props.name}</span>, Side: this.props.data.Side });
			}

			var bluePlayers = fellowPlayers.filter(function (a) {
				return a.Side === 'blue';
			});
			var purplePlayers = fellowPlayers.filter(function (a) {
				return a.Side === 'purple';
			});

			var html = bluePlayers.map(function (a, index) {
				return <tr>
							<td className="padding-right-xl padding-left-m" style={{'background-color':'#8cd6f4'}}>{bluePlayers[index].SummonerName}</td>
							<td className="padding-right-m" style={{'background-color':'#8cd6f4'}}>{bluePlayers[index].ChampionName}</td>
							<td className="padding-left-l padding-right-l"><span style={{'border-right':'1px #ccc solid'}}> </span></td>
							<td className="padding-right-xl padding-left-m" style={{'background-color':'#caa4cc'}}>{purplePlayers[index].SummonerName}</td>
							<td className="padding-right-m" style={{'background-color':'#caa4cc'}}>{purplePlayers[index].ChampionName}</td>
						</tr>;
			});

			return this.transferPropsTo(
						<div className='padding-m'>
							<table>
								<thead>
									{html}
								</thead>
							</table>
						</div>);
		}
	});

	return MatchDetails;
})