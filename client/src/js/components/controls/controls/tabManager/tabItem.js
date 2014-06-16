/** @jsx React.DOM */
define(['react'], function (React) {
	var TabItem = React.createClass({
		render: function () {
			var className = this.props.isActive ? 'active' : '';
			return this.transferPropsTo(<li className={className} onClick={this.onClick} >
											<a className="hitarea" href="#">{this.props.text}</a>
										</li>);
		},

		onClick: function () {
			this.props.onClick(this.props.index);
		}
	});

	return TabItem;
})