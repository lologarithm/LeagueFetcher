/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', '/js/components/controls/matchHistoryPanel.js'], function ($, React, moment, summonerStore, MatchHistoryPanel) {
	var MatchHistoryView = React.createClass({displayName: 'MatchHistoryView',
		getInitialState : function () {
			return {
				championInputName: '',
				championName: ''
			};
		},

		render: function () {
			var historyPanel = null;
			if(this.state.championName !== '') {
				historyPanel = MatchHistoryPanel( {name:this.state.championInputName} );
			}

			return React.DOM.div(null, 
						React.DOM.div(null, 
							React.DOM.input( {type:"text", value:this.state.championInputName, onChange:this.onChange} ),
							React.DOM.button( {type:"button", onClick:this.onClick} , "Search")
						),
						historyPanel
					);
		},

		onChange: function (e) {
			this.setState({ championInputName : e.target.value });
		},

		onClick: function () {
			this.setState({ championName : this.state.championInputName });
		}
	});

	return MatchHistoryView;
})