/** @jsx React.DOM */
define(['react'], function (React) {
	var LocalFileReader = React.createClass({

		render: function () {
			return this.transferPropsTo(
				<div>
					<input type="file" ref='file' className='padding-right-xl padding-left-xl' accept=".csv" style={{'display':'inline-block'}} />
					<button type="button" className="btn btn-small btn-primary" onClick={this.uploadBtnClick}>Upload</button>
				</div>
			);
		},

		uploadBtnClick: function() {
			var files = this.refs['file'].getDOMNode().files;

			if (!files.length) {
				alert('Please select a file!');
				return;
			}

			var file = files[0];
			var reader = new FileReader();

			// If we use onloadend, we need to check the readyState.
			reader.onloadend = $.proxy(function(evt) {
				if (evt.target.readyState == FileReader.DONE) { // DONE == 2
					var fileName = $(this.refs['file'].getDOMNode()).val().split('\\');
					fileName = fileName[fileName.length - 1];


					this.props.callback(evt.target.result, fileName);
				}
			}, this);

			var blob = file.slice(0, file.size);
			reader.readAsBinaryString(blob);
		}
	});

	return LocalFileReader;
})