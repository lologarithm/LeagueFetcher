define(['jquery', 'rx', '$rx'],

function ($, rx, $rx) {

    var Draggable = {
        //
        //  @options
        //  {
        //      parent : The container you want to bind the draggable within
        //      callback : function to be called when the mouseDrags event fires, function({ xChange, yChange })
        //  }
        //
        getInitialState: function() {
            return {
                top : 0,
                left : 0
            };
        },

        draggable: function(options) {
            options = options || {};

            var _doc = $(document), _draggable = $(this.getDOMNode()), _parent = options.parent || _draggable.parent();
            

            // get the stream of events from the mousedown, mousemove and mouseup events
            var mouseDown = _draggable.mousedownAsObservable(),
                mouseMove = _doc.mousemoveAsObservable(),
                mouseUp = _doc.mouseupAsObservable();
            
            

            // get the changes in the X and Y direction as the mouse moves
            var mouseMoves = mouseMove.skip(1).zip(mouseMove, function (left, right) {
                return {
                    xChange: left.clientX - right.clientX,
                    yChange: left.clientY - right.clientY,
                    clientX: right.clientX,
                    clientY: right.clientY
                };
            });
            
            
            // for each mouse down event, get all the subsequent changes in the clientX and
            // clientY values resulting from the mouse move events until mouse up event occurs
            var mouseDrags = mouseDown.selectMany(function (md) {
                var $target = $(md.target);

                if ($target.closest('[data-role="no-drag"]').length === 0 || $target.closest('[data-role="drag"]').length > 0) {
                    return mouseMoves.takeUntil(mouseUp);
                } else {
                    return Rx.Observable.empty();
                }
            });


            //Store the desired offset each time the mouse goes down
            // The desired offset, is the offset between the top-left corner of the element, 
            // and the mouses position. Adding this gives a better "feel" when doing dragging
            var desiredOffset = { xMouse : 0, yMouse : 0 };
            mouseDown.subscribe(function (mouseEvent) {
                desiredOffset = { xMouse : mouseEvent.clientX - _draggable.offset().left, yMouse : mouseEvent.clientY - _draggable.offset().top };
            });

            mouseDrags.subscribe($.proxy(function (mouseEvent) {
                //Compute the applied zoom throughout the layers
                var _zoom = Number(_parent.css('zoom'));
                var tempParent = _parent;
                while(tempParent[0] != $('body')[0]) {
                    _zoom *= Number(tempParent.css('zoom'));
                    tempParent = tempParent.parent();
                }

                mouseEvent.xChange /= _zoom;
                mouseEvent.yChange /= _zoom;

                //Calculate if the current offset is appropriate, 
                // if the mouse is 'out of bounds', in context to the desired, don't move the popup
                // else, move it as we would normally
                _dragOffsetX = (mouseEvent.clientX - _draggable.offset().left);
                _dragOffsetY = (mouseEvent.clientY - _draggable.offset().top);
                mouseEvent.xChange = desiredOffset.xMouse - _dragOffsetX >= 0 ? Math.min(0, mouseEvent.xChange) : Math.max(0, mouseEvent.xChange);
                mouseEvent.yChange = desiredOffset.yMouse - _dragOffsetY >= 0 ? Math.min(0, mouseEvent.yChange) : Math.max(0, mouseEvent.yChange);


                //Set up our boundries for the drag container
                var _minLeft = 0,
                _minTop = 0,
                _maxLeft = _parent.width() - _draggable.outerWidth(),
                _maxTop = _parent.height() - _draggable.outerHeight();
                
                // change the left and top, set margin to 0 as its an absolute element,
                // and we don't need any position adjustments for this
                //
                // In previous instances, having margin caused display issues 
                this.setState({
                    left: Math.min(Math.max(this.state.left + mouseEvent.xChange, _minLeft), _maxLeft),
                    top: Math.min(Math.max(this.state.top + mouseEvent.yChange, _minTop), _maxTop),
                    margin: '0px'
                });
            }, this));

            //Subscribe the specified callback, so it knows when the element 
            // is being dragged, and can act accordingly
            if(options && options.callback)
                mouseDrags.subscribe(options.callback);

            return _draggable;
        }
    };

    return Draggable;
});