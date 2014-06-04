/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', '/js/components/controls/matchHistoryPanel.js'], function ($, React, moment, summonerStore, MatchHistoryPanel) {
	var MatchHistoryView = React.createClass({
		getInitialState : function () {
			return {
				championInputName: '',
				championName: ''
			};
		},

		render: function () {
			var inputArea = <div className='flexContainer'>
								<input ref="txtInput" type="text" className='flex1' value={this.state.championInputName} onChange={this.onChange} onKeyDown={this.onKeyDown} />
								<button type="button" className='flexNone btn btn-primary' style={{'margin-left':'3px'}} onClick={this.onClick} >Search</button>
							</div>;

			if(this.state.championName !== '') {
				return <div style={{'width':'500px', 'margin':'50px auto 0'}}>
							<h3>League Fetcher</h3>
							{inputArea}
							<MatchHistoryPanel className="margin-top-m" name={this.state.championName} />
						</div>;
			} else {
				return <div style={{'width':'500px', 'margin':'200px auto 0'}}>
							<h3>League Fetcher</h3>
							{inputArea}
						</div>;
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