define([], 

function () {
	// A list of all the callbacks bound to the stores
	var _callbacks = [];


	// Class to handle all of the events in the app. It invokes callbacks to the stores to create
	//   circular data flow throughout the application.
	var Dispatcher = {

		/**
		 * Registers the callback to be called whenever an event is dispatched
		 * @callback : The function to be invoked onChange
		 *
		*/
		register: function (callback) {
			_callbacks.push(callback);

			return _callbacks.length - 1;
		},

		dispatch: function (action) {
			_callbacks.forEach(function (a){ 
				a(action);
			});
		}
	};

	return Dispatcher;
});