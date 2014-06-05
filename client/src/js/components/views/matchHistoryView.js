/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', '/js/components/controls/matchHistoryPanel.js'], function ($, React, moment, summonerStore, MatchHistoryPanel) {
	var MatchHistoryView = React.createClass({

		render: function () {
			if(this.props.name !== '') {
				return <div style={{'width':'500px', 'margin':'10px auto 0'}}>
							<h3>Match History</h3>
							<MatchHistoryPanel className="margin-top-m" name={this.props.name} />
						</div>;
			} else {
				return <div></div>;
			}
		},
	});

	return MatchHistoryView;
})