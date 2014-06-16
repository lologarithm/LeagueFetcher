/** @jsx React.DOM */
define(['react', 
'./tabItem.js', 
'./tabManagerMixin.js'], 

function (React, TabItem, tabManagerMixin) {

	var HorizontalTabManager = React.createClass({

		mixins: [tabManagerMixin],

		render: function () {
			this.props.children = this.props.children || [];
			var propsName = this.props.tabClass ? this.props.tabClass : '';

			var tabs = this.props.tabs.map($.proxy(function (item, index) {
				return <span onClick={this.itemClicked(index)} className={index === this.state.active ? 'active ' + propsName : propsName} 
																style={{'display':'inline-block', 'padding-right':'15px', 'padding-left':'15px'}}>
							{item}
						</span>;
			}, this));

			return this.transferPropsTo(<div className='flexContainer flexColumn'>
											<div className='flexNone' style={this.props.tabStyle}>
												<div className="tabManager">
													{tabs}
												</div>
											</div>
											<div className='flex1' style={{'overflow':'auto'}}>
												{this.props.children[this.state.active]}
											</div>
											<div />
										</div>);
		}
	});

	return HorizontalTabManager;
})