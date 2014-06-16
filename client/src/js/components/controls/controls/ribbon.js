define(['react'], function(React) {
    var Ribbon = React.createClass({
        render: function () {
            return this.transferPropsTo(
                <div classname="padding-right-xl padding-left-xl">
                    <div role="tabpanel" classname="tab-pane active" aria-labelledby="fileTabToggle3">
                        <div role="menubar" classname="btn-toolbar btn-toolbar-condensed">
                            {this.props.children}
                        </div>
                    </div>
                </div>);
        }
    })
});

