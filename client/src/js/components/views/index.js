/** @jsx React.DOM */
define(['jquery', 'react', 'js/components/views/matchHistoryView.js', 'js/components/views/rankedStatsView.js'], function ($, React, MatchHistoryView, RankedStatsView) {
	var IndexView = React.createClass({
		getInitialState : function () {
			return {
				championInputName: '',
				championName: ''
			};
		},

		render: function () {
			var style = this.state.championName !== '' ? {'width':'500px', 'margin':'50px auto 20px'} : {'width':'500px', 'margin':'200px auto 20px'};

			return <div className="flexContainer flexColumn">
						<div classname='flexContainer flexColumn flexNone marginAnimate' style={style}>
							<div className='flexNone margin-bottom-s' style={{"font-size": 25}}>
								Enter Summoner Name:
							</div>
							<div className='flexContainer flexNone marginAnimate'>
								<input ref="txtInput" type="text" className='flex1' value={this.state.championInputName} onChange={this.onChange} onKeyDown={this.onKeyDown} />
								<button type="button" className='flexNone btn btn-primary' style={{'margin-left':'3px'}} onClick={this.onClick} >Search</button>
							</div>
						</div>
						<div className='flexContainer flex1'>
							<MatchHistoryView className="flex1" name={this.state.championName} />
							<RankedStatsView className="flex1" name={this.state.championName} />
						</div>
					</div>
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

	function init() {
		React.renderComponent(<IndexView />, $('#container')[0]);
	}


	return init;
})