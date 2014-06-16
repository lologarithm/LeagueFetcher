/** @jsx React.DOM */
define(['react'], function (React) {
	var DropdownMenu = React.createClass({

		render: function () {

			var items = this.props.items.map($.proxy(function (item, index) {
				return <li>
				            <a className="hitarea" onClick={this.props.callbacks[index]} onMouseDown={this.preventBubbling} tabindex="-1" href="#">{item}</a>
				        </li>;
			}, this));


			return this.transferPropsTo(
				<div className="dropdown" style={{'z-index':'10000'}}>
				    <ul className="dropdown-menu" style={{'display':'block', 'overflow':'auto','max-height':this.props.menuHeight}}>
				        {items}
				    </ul>
				</div>
			);
		},

		preventBubbling: function(e) {
			e.stopPropagation();
		}
	});

	return DropdownMenu;
})