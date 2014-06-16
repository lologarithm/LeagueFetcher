/** @jsx React.DOM */
define(['react'], function (React) {
	var ListAccordian = React.createClass({

		getInitialState: function(){
            return {
                isExpanded : true
            };
        },

        componentDidMount: function () {
			if(this.props.isExpanded !== undefined) {
				this.setState({isExpanded : this.props.isExpanded});
			}
        },

		render: function () {
			var iconClass = 'icon ' + (this.state.isExpanded ? 'icon-chevron-down' : 'icon-chevron-right');

			return this.transferPropsTo(<div className="margin-bottom-l">
											<div className="padding-bottom-s margin-bottom-m" style={{'border-bottom' : '1px #000 solid', 'font-size':'12px'}}>
												<i className={iconClass} onClick={this._onClick} style={{'font-size':'12px'}}/>
												{this.props.title}
											</div>
											<div style={{'display' : (this.state.isExpanded ? 'block' : 'none')}}>
												{this.props.children}
											</div>
										</div>);
		},

		_onClick: function (e) {
			this.setState({isExpanded : !this.state.isExpanded});
			e.stopPropagation();
		}
	});

	return ListAccordian;
})