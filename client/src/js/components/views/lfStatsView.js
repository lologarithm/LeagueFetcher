/** @jsx React.DOM */
define(['jquery', 'react', 'moment', 'js/stores/summonerStore.js', 'js/components/controls/rankedStatsElement.js'], 

function ($, React, moment, summonerStore, RankedStatsElement) {
	var RankedStatsView = React.createClass({
		getInitialState: function () {
			return {
				dataUpdate: false
			};
		},
		
		render: function () {
			if (this.props.name !== '') {
				return this.transferPropsTo(
					<div>
						{"League Fetcher's stats module is coming! Keep an eye out in the coming weeks!"}
					</div>
				);
			} else {
				return <div></div>;
			}
		},

		componentDidMount: function ( ){
			summonerStore.addChangeListener(this.onChange);
		},

		componentWillUnmount: function () {
			summonerStore.removeChangeListener(this.onChange);
		},

		componentDidUpdate: function () {
			this.state.dataUpdate = false;
		},

		shouldComponentUpdate: function (nextProps) {
			return this.props.name !== nextProps.name || this.state.dataUpdate;
		},

		onChange: function () {
			this.setState({ dataUpdate : true });
		}
	});

	return RankedStatsView;
})