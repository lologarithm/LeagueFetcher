/** @jsx React.DOM */
define(['react'], function (React) {
	var Modal = React.createClass({

		getInitialState: function(){
            return {
                isVisible : false,
            };
        },

        componentDidMount: function () {
			this.setState( { isVisible : this.props.isVisible });
        },

        componentWillReceiveProps: function(nextProps) {
			this.setState( { isVisible : nextProps.isVisible });
        },

		render: function () {
			return this.transferPropsTo(<div className="modal fade in" style={{'display' : this.state.isVisible ? 'block' : 'none' }} >
						<div className="modal-dialog">
							<div className="modal-content">
								<div className="modal-header">
					                <h4 className="modal-title" id="myModalLabel">{this.props.header}
					                	<button type="button" className="close" onClick={this.closeModal} data-dismiss="modal" aria-hidden="true" style={{'position':'relative','padding':'0 4px 0 15px','line-height':'inherit','text-align':'middle'}}>&times;</button>
					                	{this.props.customHeader}
				                	</h4>
					            </div>
					            <div className="modal-body margin-left-xl margin-right-xl">
					            	{this.props.children}
					            </div>
					            <div className="modal-actions modal-footer">
					                {this.props.footer}
					            </div>
							</div>
						</div>
					</div>);
		},

		closeModal: function () {
			this.setState({isVisible : !this.state.isVisible});

			if(this.props.onClose)
				this.props.onClose();
		}
	});

	return Modal;
})