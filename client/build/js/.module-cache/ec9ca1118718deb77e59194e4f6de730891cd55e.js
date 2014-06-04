/** @jsx React.DOM */
define(['jquery', 'react'], function ($, React) {
	var IndexView = React.createClass({displayName: 'IndexView',
		render: function () {
			return React.DOM.div(null, "TEST");
		}
	});


	React.renderComponent(IndexView(), $('#body')[0]);

	return true;
})