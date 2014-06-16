define(['react'], function (React){
	var obj = {
		propTypes: {
			tabs: React.PropTypes.array.isRequired, // Array of react objects to use as tab items. They are wrapped in spans for the list
			active : React.PropTypes.number, // The index of the active tab. NOT REQUIRED, but can be passed to override the controls state
			onActiveChanged : React.PropTypes.func //Callback with the index of hte new tab switched to
		},
		
		getInitialState: function(){
            return {
                active : 0,
            };
        },

        componentDidMount: function () {
        	if(this.props.active)
				this.setState({active : this.props.active});
        },

        componentWillReceiveProps: function(nextProps) {
			if(nextProps.active != this.state.active && nextProps.active !== undefined)
				this.setState({active : nextProps.active});
        },

		itemClicked: function (index) {
			return $.proxy(function () {
				this.setState({active : index});

				if(this.props.onActiveChanged) {
					this.props.onActiveChanged(index);	
				}
			}, this);
		}
	};

	return obj;
});