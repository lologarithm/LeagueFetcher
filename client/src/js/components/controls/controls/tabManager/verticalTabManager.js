/** @jsx React.DOM */
define(['react', 
'./tabItem.js', 
'./tabManagerMixin.js'], 

function (React, TabItem, tabManagerMixin) {
	var VerticalTabManager = React.createClass({

		mixins: [tabManagerMixin],

		render: function () {
			this.props.children = this.props.children || [];

			var tabs = this.props.tabs.map($.proxy(function (tab, index) {
				return this.state.active === index ? <div key={index} onClick={this.itemClicked(index)} className="selected">{tab}</div> : <div key={index} onClick={this.itemClicked(index)}>{tab}</div>;
			}), this);

			return this.transferPropsTo(<div className="verticalTabManager flexContainer">											
											<div className='tabContents flex1 scroll flexContainer' style={{'overflow-x':'hidden'}}>
												{this.props.children[this.state.active]}
											</div>
											<div className='tabItems flexNone flexContainer flexColumn'>
												{tabs}
												<div className="flex1"></div>
											</div>
										</div>);
		}
	});

	return VerticalTabManager;
})