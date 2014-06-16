/** @jsx React.DOM */
define(['react', './dropdown.js'], 

function (React, Dropdown) {
	var AutoCompleteInput = React.createClass({

		getInitialState: function(){
            return {
                value: this.props.value,
                displayDropdown : false
            };
        },

        componentWillReceiveProps: function(nextProps) {
			this.setState( { value : nextProps.value });
        },

		render: function () {
			this.props.options = this.props.options || [];

			var controls = this.props.options.map($.proxy(function (a) {
				return <span>
							{a.option}
						</span>;
			}, this));

			var callbacks = this.props.options.map($.proxy(function (a) {
				return $.proxy(function () {
					this.onChange({ target : { value : a.str } });
					this.setState({ displayDropdown : false });

					var $field = $('input', this.getDOMNode());
					var oldVal = $field.val();

					$field.focus().val('').val(oldVal);
				}, this);
			}, this));

			return <div>
						{ this.transferPropsTo(<input ref="input" type="text" value={this.state.value} onChange={this.onChange} />) }
						<Dropdown items={controls} menuHeight="200px" callbacks={callbacks} style={{'display': (!this.state.displayDropdown || this.state.value == '' || this.props.options.length == 0 ? 'none':'block')}} />
					</div>;
		},

		onChange: function (evt) {
			this.state.value = evt.target.value;
			this.state.displayDropdown = true;
        	this.setState(this.state);

        	if(this.props.onChange) {
        		this.props.onChange(evt);
        	}
		}
	});

	return AutoCompleteInput;
})