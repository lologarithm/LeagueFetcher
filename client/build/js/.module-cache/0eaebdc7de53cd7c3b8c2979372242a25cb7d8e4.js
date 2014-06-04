/** @jsx React.DOM */
define(['jquery', 'react', 'js/components/control/matchHistoryPanel.js'], function ($, React, MatchHistoryPanel) {
	var IndexView = React.createClass({displayName: 'IndexView',
		render: function () {
			return MatchHistoryPanel(null);
		}
	});

	function init() {
		React.renderComponent(IndexView(null ), $('#container')[0]);
	}


	

	return init;
})