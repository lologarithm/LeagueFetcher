/** @jsx React.DOM */
define(['jquery', 'react', 'js/components/views/matchHistoryView.js'], function ($, React, MatchHistoryView) {
	var IndexView = React.createClass({
		render: function () {
			return <MatchHistoryView />;
		}
	});

	function init() {
		React.renderComponent(<IndexView />, $('#container')[0]);
	}


	return init;
})