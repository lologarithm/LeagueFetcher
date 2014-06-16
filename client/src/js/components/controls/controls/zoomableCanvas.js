/** @jsx React.DOM */
define(['react', 'jquery'], function (React, $) {
	var ZoomableCanvas = React.createClass({

		MAX_ZOOM : 2.5,
		MIN_ZOOM : .4,

		getInitialState: function(){
            return {
                zoom : 1
            };
        },

        componentDidMount: function () {
			$(window).bind('mousewheel', $.proxy(function (e) {
				var $this = $(this.getDOMNode());
				var offset = $(this.getDOMNode()).offset();
				var x = event.clientX;
				var y = event.clientY;

				if(x > offset.left && x < offset.left + ($this.outerWidth(true) * this.state.zoom) && y > offset.top && y < offset.top + ($this.outerHeight(true) * this.state.zoom)) {
					var delta = e.originalEvent.wheelDelta;
					var newZoom = 1 + ((this.props.sensitivty || .01) * delta / 120);

					this.setState({ zoom : Math.min(this.MAX_ZOOM, Math.max(this.state.zoom * newZoom, this.MIN_ZOOM)) });

					if(this.props.zoomCallback)
						this.props.zoomCallback(this.state.zoom);
				}
			}, this));
        },

		render: function () {
			var iconClass = 'icon ' + (this.state.isExpanded ? 'icon-chevron-down' : 'icon-chevron-right');

			return this.transferPropsTo(<div style={{'zoom': this.state.zoom * 100 + '%'}}>
											{this.props.children}
										</div>);
		}
	});

	return ZoomableCanvas;
})