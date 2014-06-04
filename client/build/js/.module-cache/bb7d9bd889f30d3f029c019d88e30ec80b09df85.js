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
			var inputArea = React.DOM.div( {className:"flexContainer"}, 
								React.DOM.input( {ref:"txtInput", type:"text", className:"flex1", value:this.state.championInputName, onChange:this.onChange, onKeyDown:this.onKeyDown} ),
								React.DOM.button( {type:"button", className:"flexNone btn btn-primary", style:{'margin-left':'3px'}, onClick:this.onClick} , "Search")
							);

			if(this.state.championName !== '') {
				return React.DOM.div( {style:{'width':'400px', 'margin':'50px auto 0'}}, 
							React.DOM.h3(null, "League Fetcher"),
							inputArea,
							MatchHistoryPanel( {className:"margin-top-m", name:this.state.championName} )
						);
			} else {
				return React.DOM.div( {style:{'width':'400px', 'margin':'200px auto 0'}}, 
							React.DOM.h3(null, "League Fetcher"),
							inputArea
						);
			}
		},

		onChange: function (e) {
			this.setState({ championInputName : e.target.value });
		},

		onKeyDown: function (e) {
			if(e.keyCode === 13) {
				this.commitSearch();
			}
		},

		onClick: function () {
			this.commitSearch();
		},

		commitSearch: function(name) {
			this.setState({ championName : this.state.championInputName });
			$(this.refs.txtInput.getDOMNode()).blur();
		}
	});

	return MatchHistoryView;
})