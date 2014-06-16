/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', '/js/components/controls/matchHistoryPanel.js'], function ($, React, moment, summonerStore, MatchHistoryPanel) {
	var MatchHistoryView = React.createClass({

		render: function () {
			if(this.props.name !== '') {
				return this.transferPropsTo(
						<div>
							<h3>Match History</h3>
							<MatchHistoryPanel style={{'width':'50%'}} className="margin-top-m" name={this.props.name} searchName={this.props.searchName} />
						</div>);
			} else {
				return <div></div>;
			}
		},
	});

	return MatchHistoryView;
})