/** @jsx React.DOM */
define(['react'], function(React) {
	var Popover = React.createClass({
		propTypes: {
			arrowPos: React.PropTypes.str,
			confirmClick : React.PropTypes.func, // The index of the active tab. NOT REQUIRED, but can be passed to override the controls state
			cancelClick : React.PropTypes.func //Callback with the index of hte new tab switched to
		},

		render: function () {
			var cancelBtn = this.props.cancelClick !== undefined ? <img src="images/formula-controls-cancel.svg" onClick={this.cancelClick} height="16" /> : '';
			var confirmBtn = this.props.confirmClick !== undefined ? <img src="images/formula-controls-confirm.svg" onClick={this.confirmClick} height="16" /> : '';

			var arrowStyle = this.props.arrowStyle || {};
			if(this.props.style['border-color']) {
				if(this.props.className.indexOf('left') !== -1) {
					arrowStyle['border-left-color'] = this.props.style['border-color'];
				} else if(this.props.className.indexOf('right') !== -1) {
					arrowStyle['border-right-color'] = this.props.style['border-color'];
				} else if(this.props.className.indexOf('top') !== -1) {
					arrowStyle['border-top-color'] = this.props.style['border-color'];
				} else if(this.props.className.indexOf('bottom') !== -1) {
					arrowStyle['border-bottom-color'] = this.props.style['border-color'];
				}
				
			}


			return this.transferPropsTo(
				<div className="popover">
				    <div className="arrow" style={arrowStyle} ></div>
					<div className="inner">
						<div className="content">
				            {this.props.children}
					    </div>	    
					    <div className="controls">
					    	{cancelBtn}
					    	{confirmBtn}
					    </div>
				    </div>
				</div>
			);
		},

		cancelClick: function (e) {
			if(this.props.cancelClick) {
				this.props.cancelClick();
			}

			e.stopPropagation();
		},

		confirmClick: function (e) {
			if(this.props.confirmClick) {
				this.props.confirmClick();
			}

			e.stopPropagation();
		}
	});

	return Popover;
});

