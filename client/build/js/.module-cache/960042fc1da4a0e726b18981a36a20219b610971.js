/** @jsx React.DOM */
define(['jquery', 'react', 'js/components/control/matchHistoryPanel.js'], function ($, React) {
	var IndexView = React.createClass({displayName: 'IndexView',
		render: function () {
			return React.DOM.div(null, "TEST");
		}
	});

	function init() {
		React.renderComponent(IndexView(null ), $('#container')[0]);
	}


	

	return init;
})