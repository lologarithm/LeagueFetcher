/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js'], function ($, React, moment, summonerStore) {
	var Popover = React.createClass({
		render: function () {
			return this.transferPropsTo(
					<div className="popover fade right in" style={{'display': 'block', 'left': '237px'}}>
						<div className="arrow" style={{'top':'55px'}}>
						</div>
						<h3 className="popover-title">
							{this.props.title}
						</h3>
						<div className="popover-content">
							{this.props.children}
						</div>
					</div>);
		}
	});

	return Popover;
})