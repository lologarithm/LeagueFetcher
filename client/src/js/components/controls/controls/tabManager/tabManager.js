/** @jsx React.DOM */
define(['react', 
'./tabItem.js', 
'./tabManagerMixin.js'], 

function (React, TabItem, tabManagerMixin) {
	var TabManager = React.createClass({

		mixins: [tabManagerMixin],

		render: function () {
			var tabItems = this.props.items.map($.proxy(function (item, index) {
				return <TabItem index={index} isActive={this.state.active === index} text={item} onClick={this.itemClicked} />;
			}, this));

			return this.transferPropsTo(<ul className="nav nav-pills nav-justified">
						{tabItems}
					</ul>);
		}
	});

	return TabManager;
})