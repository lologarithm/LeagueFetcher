/** @jsx React.DOM */
define(['react'], function (React) {
	var CollapsingPanel = React.createClass({
		getInitialState: function () {
			return {
				isHidden: false,
			};
		},

		render: function () {
			var className = 'heightWidthAnimation flexContainer flexColumn';

			var iconClass = 'icon right hitarea icon-md ';
			iconClass += this.state.isHidden ? ' icon-backward' : ' icon-forward';
 	
 			if(this.state.isHidden) {
 				return this.transferPropsTo(<div className={className} style={{'background-color':'#f0f0f0'}}>
												<div className="padding-left-m padding-right-s padding-bottom-s padding-top-s">
													<i onClick={this.collapseClick} className={iconClass} style={{'color':'#aaa'}} />
												</div>
											</div>);
 			} else {
 				return this.transferPropsTo(<div className={className} style={{'background-color':'#f0f0f0'}}>
												<div className="flexNone flexContainer" >
													<div className="flex1">{this.props.header}</div>
													<i className="flexNone" onClick={this.collapseClick} className={iconClass} style={{'color':'#aaa'}} />
												</div>
												<div className="flex1 flexContainer">
													{this.props.children}
												</div>
											</div>);
 			}
			
		},

		collapseClick: function () {
			this.setState({isHidden : !this.state.isHidden});

			if(this.props.callback) {
				this.props.callback(this.state.isHidden);
			}
		}
	});

	return CollapsingPanel;
})