/** @jsx React.DOM */
define([
'jquery', 
'react', 
'js/components/views/generalView.js',
'js/components/views/matchHistoryView.js',
'js/components/views/rankedStatsView.js',
'js/components/views/lfStatsView.js',
'js/components/controls/controls/tabManager/verticalTabManager.js'], 

function ($, React, GeneralView, MatchHistoryView, RankedStatsView, LFStatsView, TabManager) {
	var IndexView = React.createClass({
		getInitialState : function () {
			return {
				championInputName: '',
				championName: ''
			};
		},

		render: function () {
			var style = this.state.championName !== '' ? {'margin': '0px auto', 'width':'125px'} : {'width':'500px', 'margin': '0px auto'};

			return <div className="flexContainer flexColumn">
						<div className='flexContainer flexColumn flexNone marginAnimate' style={{'width':'500px', 'margin':'50px auto 20px'}}>
							<img className="widthAnimate" src="/imgs/logo.png" style={style} />
							<div className='flexNone margin-bottom-s' style={{"font-size": 25}}>
								Summoner Search
							</div>
							<div className='flexContainer flexNone'>
								<input ref="txtInput" type="text" className='flex1' value={this.state.championInputName} onChange={this.onChange} onKeyDown={this.onKeyDown} />
								<button type="button" className='flexNone btn btn-primary' style={{'margin-left':'3px'}} onClick={this.onClick} >Search</button>
							</div>
						</div>
						<TabManager style={{'display':(this.state.championName !== '' ? '':'none'),'width':'95%','max-width':'1200px', 'margin':'0 auto 20px', 'background-color':'#e9eaed'}} 
									className='left' 
									tabs={['General', 'Match History', 'Ranked Stats', 'LF Stats']}>

							<GeneralView className="flex1" name={this.state.championName} searchName={this.commitSearch} />
							<MatchHistoryView className="flex1" name={this.state.championName} searchName={this.commitSearch} />
							<RankedStatsView className="flex1" name={this.state.championName} searchName={this.commitSearch} />
							<LFStatsView className="flex1" name={this.state.championName} searchName={this.commitSearch} />
						</TabManager>
					</div>
		},

		onChange: function (e) {
			this.setState({ championInputName : e.target.value });
		},

		onKeyDown: function (e) {
			if(e.keyCode === 13) {
				this.commitSearch(this.state.championInputName);
			}
		},

		onClick: function () {
			this.commitSearch(this.state.championInputName);
		},

		commitSearch: function(name) {
			this.setState({ championName : name, championInputName: name });
			$(this.refs.txtInput.getDOMNode()).blur();
		}
	});

	function init() {
		React.renderComponent(<IndexView />, $('#container')[0]);
	}


	return init;
})