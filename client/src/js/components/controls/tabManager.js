/** @jsx React.DOM */
define(['jquery', 'react'], function ($, React) {
	var Popover = React.createClass({
		render: function () {
			var tabs = this.props.tabs.map(function (a ){ 
				return <div>{a}</div>;
			});

			return this.transferPropsTo(
					<div className="flexContainer">
						<div style={{'border-right':'#ccc solid 1px'}}>
							{tabs}
						</div>
						<div className="flex1">
							this.props.children[this.state.activeTab]
						</div>
					</div>);
		}
	});

	return Popover;
})