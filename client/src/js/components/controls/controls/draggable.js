/** @jsx React.DOM */
define(['react'], function (React) {
	var Draggable = React.createClass({
		mixins : [DraggableMixin],

		componentDidMount: function () {
			this.draggable({ callback : this.props.dragCallback });
        },

		render: function () {
			return this.transferPropsTo(<div>
											{this.props.children}
										</div>);
		}
	});

	return Draggable;
})