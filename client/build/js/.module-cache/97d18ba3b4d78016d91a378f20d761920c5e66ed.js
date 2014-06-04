define(['jquery', 'react'], function ($, React) {
	var IndexView = React.createClass({
		render: function () {
			return <div>TEST</div>;
		}
	});


	React.renderComponent(IndexView(), $('#body')[0]);
})