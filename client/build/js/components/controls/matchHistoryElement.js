/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js'], function ($, React, moment, summonerStore) {
	var MatchHistoryPanel = React.createClass({displayName: 'MatchHistoryPanel',
		render: function () {
			var kda = (this.props.data.Stats.ChampionsKilled + this.props.data.Stats.Assists) / (this.props.data.Stats.NumDeaths || 1);
			var elmStyle = this.props.data.Stats.Win ? {'border-top':'2px solid green'} : {'border-top':'2px solid red'};

			return React.DOM.div( {className:"matchHistoryElement padding-m margin-bottom-m", style:elmStyle} , 
						React.DOM.div(null, 
							React.DOM.span(null, "Champion Played: ", this.props.data.ChampionName),
							React.DOM.span( {style:{'display':'inline-block', 'float':'right'}}, 
								React.DOM.span( {style:{'color':'green'}}, this.props.data.Stats.ChampionsKilled),"/",
								React.DOM.span( {style:{'color':'red'}}, this.props.data.Stats.NumDeaths),"/",
								React.DOM.span( {style:{'color':'orange'}}, this.props.data.Stats.Assists)
							)
						),
						React.DOM.div(null, 
							React.DOM.span(null, "Date Played: ",  moment(this.props.data.CreateDate).format('MMMM Do YYYY, h:mm a') ),
							React.DOM.span( {style:{'display':'inline-block', 'float':'right'}}, "KDA: ", kda)
						),
						React.DOM.div(null, "Game Mode: ",  this.props.data.GameMode ),
						React.DOM.div(null, "IP Earned: ",  this.props.data.IpEarned )
					);
		}
	});

	return MatchHistoryPanel;
})