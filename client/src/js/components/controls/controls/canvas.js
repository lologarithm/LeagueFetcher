define(['react', '$'],

function (React, $) {

	var Canvas = React.createClass({
		render: function () {
			return <Canvas ref="canvas" />;
		},

		getInitialState: function () {
			return {
				drawFunctions : {}
			};
		},

		componentDidMount: function () {
			this.redraw(true);
		},

		redraw: function (clearCanvas) {
			var $canvas = $(this.getDOMNode());
			var ctx = $canvas[0].getContext('2d');

			if(clearCanvas === undefined || clearCanvas)
				ctx.clearRect(0, 0, $canvas.outerWidth(), $canvas.outerHeight());

			for(var callback in this.canvas.drawFunctions) {
				this.canvas.drawFunctions[callback](this.canvas.$canvas);
			}
		},

		registerDrawCall: function (key, drawFunc) {
			this.canvas.drawFunctions[key] = drawFunc;
		},

		removeDrawCall: function (key) {
			delete this.canvas.drawFunctions[key];
		}
	});

	return Canvas;
});