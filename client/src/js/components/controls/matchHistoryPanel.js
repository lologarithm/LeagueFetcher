/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', 'js/stores/matchStore.js', './matchHistoryElement.js', './matchDetails.js', './popover.js'], 

function ($, React, moment, summonerStore, matchStore, MatchHistoryElement, MatchDetails, Popover) {
	var MatchHistoryPanel = React.createClass({
		getInitialState: function () {
			return {
				target: null,
				targetPosition: {top:0,left:0}
			}
		},

		render: function () {
			if(summonerStore.getMatchHistory(this.props.name) && summonerStore.getMatchHistory(this.props.name).Games && summonerStore.getMatchHistory(this.props.name).Games.length > 0) {
				var popover = summonerStore.getMatchHistory(this.props.name).Games.filter($.proxy(function (a) {
					return this.state.target == a.GameId;
				}, this)).map($.proxy(function (a) {
					return <Popover title={a.GameId} style={{'top': this.state.targetPosition.top, 'left': this.state.targetPosition.left, 'max-width':'100000px'}} >
								<MatchDetails data={matchStore.getMatch(a.GameId, summonerStore.getMatchHistory(this.props.name).SummonerId, a, this.props.name)} searchName={this.props.searchName} name={this.props.name} />
							</Popover>;
				}, this))[0];

				var items = summonerStore.getMatchHistory(this.props.name).Games.map($.proxy(function (a) {
					return <MatchHistoryElement data-id={a.GameId} onClick={this.elementClicked} data={a} />;
				}, this));

				return this.transferPropsTo(
								<div>
									{popover}
									{items}
								</div>);
			} else {
				return this.transferPropsTo(<div>No recent games found</div>);
			}
			
		},

		componentDidMount: function ( ){
			summonerStore.addChangeListener(this.onChange);
			matchStore.addChangeListener(this.onChange);

			//Bind Closing the popover
			$(document).bind('mousedown', this.mouseDown);
		},

		componentWillUnmount: function () {
			summonerStore.removeChangeListener(this.onChange);
			matchStore.removeChangeListener(this.onChange);

			$(document).unbind('mousedown', this.mouseDown);
		},

		mouseDown: function (e) {
			if($(e.target).closest($(this.getDOMNode())).length === 0) {
				this.setState({ target: null });
			}
		},

		onChange: function () {
			this.setState({});
		},

		elementClicked: function (e) {
			var $target = $(e.target).closest('[data-id]');

			this.setState({ target : $target.attr('data-id'), targetPosition: { top: $target.offset().top, left: $target.offset().left + $target.outerWidth(true) } });
		}
	});

	return MatchHistoryPanel;
});