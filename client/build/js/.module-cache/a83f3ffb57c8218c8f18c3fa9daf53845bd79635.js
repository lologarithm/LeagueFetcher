/** @jsx React.DOM */
define(['jquery', 'react', 'js/components/controls/matchHistoryView.js'], function ($, React, MatchHistoryView) {
	var IndexView = React.createClass({displayName: 'IndexView',
		render: function () {
			return MatchHistoryView(null);
		}
	});

	function init() {
		React.renderComponent(IndexView(null ), $('#container')[0]);
	}


	

	return init;
})