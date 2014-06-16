/** @jsx React.DOM */
define(['jquery'], function ($) {
	return {
		getInitialState: function () {
			return {
				contextMenu: {
					top: 0,
					left: 0,
					target: { col:0, row:0 },
					isDisplaying: false
				}
			};
		},

		componentDidMount: function () {
			//Bind context menu
			$this = $(this.getDOMNode());
			$this[0].oncontextmenu = function() {return false;};

			$(document).bind('mousedown', this.mouseDown);
		},

		componentWillUnmount: function () {
			$(document).unbind('mousedown', this.mouseDown);
		},

		mouseDown: function (e) {
			if($(e.target).closest($(this.getDOMNode())).length === 0) {
				this.state.contextMenu.isDisplaying = false;
				this.setState(this.state);
			}
		},

		openContext: function (e) {
			$this = $(this.getDOMNode());
			e.stopPropagation();

			if(e.button == 2) { 
				this.state.contextMenu.isDisplaying = true;
				this.state.contextMenu.left = $this.scrollLeft() + e.clientX - $this.offset().left;
				this.state.contextMenu.top = $this.scrollTop() + e.clientY - $this.offset().top;

				if(this.contextOpened)
					this.contextOpened(e);

				this.setState(this.state);

				return false; 
			} else {
				if(this.state.contextMenu.isDisplaying) {
					this.state.contextMenu.isDisplaying = false;

					if(this.contextClosed)
						this.contextClosed(e);

					this.setState(this.state);
				}

				return true;
			}
		}
	};
});