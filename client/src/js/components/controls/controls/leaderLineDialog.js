define(['jquery', 'lodash', 'rx', '$rx', 
	'text!templates/components/spreadsheet/spreadsheet.detail.html', 'text!templates/components/spreadsheet/spreadsheet.detail.header_edit.html', 'text!templates/components/spreadsheet/spreadsheet.detail.cell_edit.html',
	'modal', 'js/components/spreadsheet/spreadsheet.canvas', 'js/components/spreadsheet/spreadsheet.multi_part_view.js'], 


function($, _, rx, $rx, template, headerTemplate, cellTemplate, modal, canvas, multiPart) 
{

	function init(options) {

		// //Make sure the required modules are installed
		if(!this.canvas)
			this.addModule(canvas, 'canvas', options);

		if(!this.multi_part) {
			this.addModule(multiPart, 'multi_part');
		}

		this.canvas.registerDrawCall('drawCanvasLines', this.detail.drawCanvasLines, false);
		this.onColumnsAddedCallbacks.push(this.detail.selectFirstHeader);
		this.data.data.observables.push(this.canvas.redraw);

		//Initialize our modal, and inject it into the container
		$(template).appendTo(options.targetContainer);

		this.detail._modal = $('#detailModal').modal({ show : false, backdrop : false });
		this.detail._modal.on('hide.wdesk.modal', $.proxy(onModalHide, this));
		this.detail._modalDialog = $('.modal-dialog', this.detail._modal);
		this.detail._modalDialog.draggable({ callback : this.detail.onDrag, parent : options.targetContainer });
		this.detail._modalDialog.append(headerTemplate);
		this.detail._modalDialog.append(cellTemplate);
		this.detail._drawOpacity = 1;


		//SETUP events for the header details dialog
		var modDlg = this.detail._modalDialog;
		$('#saveHeaderBtn', modDlg).click(this.detail.onModalSaveHeader);
		$('[data-id="dropdown"]', modDlg).dropdown();
		$('li', modDlg).click(function () {
			$('[data-id="dropdown_selection"]', modDlg).html($('a', $(this)).html());
		});
		this.multi_part.bindSelectionElement(modDlg);


		//SETUP events for the cell details dialog
		$('#saveCellBtn', modDlg).click(this.detail.onModalSaveCell);

		this.viewContainer.on('dblclick', '.header > td', $.proxy(onDoubleClick('header'), this));
		this.viewContainer.on('dblclick', 'td[data-cell-location]', $.proxy(onDoubleClick('cell'), this));
	}

	function onDrag() {
		this.canvas.redraw();
	}

	function onDoubleClick(type) {
		return function (e) {
			var target = $(e.currentTarget);
			this.detail.launchAtTarget(target, type);
		}	
	}

	function onModalHide(e) {
		//To match the fade on the modal, we have to fade the canvas as well
		// 17 millisecond delay = 60 FPS
		// 150 is the current fade time for the modal
		var interval = window.setInterval($.proxy(function () {
			this.detail._drawOpacity -= .112;
			this.canvas.redraw();
		}, this), 17);

		setTimeout($.proxy(function () {
			clearInterval(interval);

			this.detail._drawOpacity = 1;
			this.detail.target = null;
			this.canvas.redraw();
		}, this), 150);
	}

	function onModalSaveHeader(e) {
		var newValue = $('#dimensionTxt').val();
		this.saveData(this.detail.target.attr('data-header-location'), newValue);		
		$('[data-property="' + this.detail.target.html() + '"]', this.viewContainer).attr('data-property', newValue);
		this.detail.target.html(newValue);

		this.detail._modal.modal('hide');
	}

	function onModalSaveCell(e) {
		var input = $('input', this.detail.target)[0];
		this.saveData(input.id, $('#attributeTxt').val());
		this.compute_all();
		this.detail._modal.modal('hide');
	}

	function launchAtTarget(target, type) {
		var offset = target.offset();

		var value = type === 'header' ? target.html() : $('input', target).val();

		this.detail.target = target;
		this.detail._modalDialog.css({ left : target.outerWidth() + offset.left + 100, top : target.outerHeight() + offset.top + 20 });

		$('[data-role="cell"], [data-role="header"]', this.detail._modalDialog).hide();
		$('[data-role="' + type + '"]', this.detail._modalDialog).show();
		$.bind_data(this.detail._modalDialog, { value : value });

		$('[data-dimension]', this.detail._modalDialog).attr('data-dimension', target.html());

		this.detail._modal.modal('show');
		this.canvas.redraw();	
	}

	function drawCanvasLines($canvas) {
		if(this.detail.target && !$.contains(document.documentElement, this.detail.target[0])) {
			this.detail.target = $('.header > [data-property="' + this.detail.target.attr('data-property') + '"]');
		}

		if(this.detail.target) {
			var ctx = $canvas[0].getContext("2d");

			//Calculate the padding on the dialog
			var horizPadding = (this.detail._modalDialog.innerWidth() - this.detail._modalDialog.width()) / 2 + 1;
			var vertPadding = (this.detail._modalDialog.innerHeight() - this.detail._modalDialog.height()) / 2 + 1;

			//Calculate the offsets of all the elements
			var targetOffsetUnadjusted = this.detail.target.offset();
			var dialogOffsetUnadjusted = this.detail._modalDialog.offset();
			var canvasOffset = $canvas.offset();
			var popupOffset = { left : dialogOffsetUnadjusted.left + horizPadding - canvasOffset.left, top : dialogOffsetUnadjusted.top + vertPadding- canvasOffset.top };
			var targetOffset = { top: targetOffsetUnadjusted.top - canvasOffset.top, left: targetOffsetUnadjusted.left - canvasOffset.left };

			
			// Green rectangle
			ctx.beginPath();
			ctx.lineWidth="3";
			ctx.strokeStyle="rgba(123, 206, 5, " + this.detail._drawOpacity + ")"; // "#7bce05" with applied opacity
			ctx.rect(targetOffset.left, targetOffset.top, this.detail.target.outerWidth(), this.detail.target.outerHeight());
			ctx.stroke();


			// Leader line
			var desiredSlope = popupOffset.top > targetOffset.top ? -3 : 3; //the desired angle of the leader
			var MIN_BUFFER = 0;											//the buffer area on the sides of the cell

			//Calculate basic things, and initialize variables for drawing the leader line
			var desiredMinLineEnd = { left : targetOffset.left + this.detail.target.outerWidth() + MIN_BUFFER, top: targetOffset.top };
			var xHitPoint = 0;
			var yHitPoint = targetOffset.top;

			var drawHorizontal = true;

			// Here is a 'zoning' of what the following if/else block covers
			// Each zone represents a different display dynamic, and is covered accordingly
			// The very middle is the cell, and the popup being there covers its contents, 
			// so we don't really care about the leader line there
			//
			//	1: Right buffer zone
			//	2: Directly Below cell
			//	3. Above, and left of the right buffer
			//	4. Left buffer zone
			//	5. Lower left region
			//	6. Everything right of the right buffer
			//
			//
			//						   | |
			//			3			   | |
			//						   | |
			// __________________ _ _ _| |
			//                 | |_ _ _|1|
			//                 | | 	   | |		6
			//				   | |     | |
			//			5	   |4|  2  | |
			//				   | |     | |
			//				   | |     | |         

			// IF the popup lies within the right-hand buffer zone, no longer require a slope, but rather just set the x-intercept
			// half way between the edge of the cell, and the popup
			if (popupOffset.left < desiredMinLineEnd.left + MIN_BUFFER && popupOffset.left > desiredMinLineEnd.left - MIN_BUFFER) {
				xHitPoint = targetOffset.left + this.detail.target.outerWidth() + (popupOffset.left - (targetOffset.left + this.detail.target.outerWidth())) / 2;
			} 

			// ELSE IF the popup is underneath the cell, and inbetween its right and left edges, simply draw a line 
			// straight up, and disabled the horizontal
			else if (targetOffset.left + this.detail.target.outerWidth() >= popupOffset.left && popupOffset.left > targetOffset.left && popupOffset.top >= targetOffset.top) {
				xHitPoint = popupOffset.left;
				yHitPoint = targetOffset.top + this.detail.target.outerHeight();
				drawHorizontal = false;
			} 

			// ELSE IF the popup is on top, and behind the right edge
			// simply draw a line down to the cell's y-coord directly, and fill in withh the horizontal
			else if (targetOffset.left + this.detail.target.outerWidth() >= popupOffset.left && popupOffset.top < targetOffset.top) {
				xHitPoint = popupOffset.left;
			} 

			// ELSE IF the popup is in the left hand buffer,
			// set the x-intercept to half way between the cell and the popup
			else if (popupOffset.left <= targetOffset.left && popupOffset.left > targetOffset.left - 2 * MIN_BUFFER) {
				xHitPoint = targetOffset.left - (targetOffset.left - popupOffset.left) / 2;
			} 

			// ELSE IF the popup is far left of the cell and beneath it, 
			// do a normal slope calculation, with a boundry on the left buffer area
			else if (popupOffset.left <= targetOffset.left - 2 * MIN_BUFFER && popupOffset.top > targetOffset.top) {
				xHitPoint = Math.min(targetOffset.left - MIN_BUFFER, (targetOffset.top - popupOffset.top) / desiredSlope + popupOffset.left);
			} 

			// ELSE do a normal slope calculation, with a boundry on the right buffer area
			else {
				xHitPoint = Math.max(desiredMinLineEnd.left, (targetOffset.top - popupOffset.top) / (-1 * desiredSlope) + popupOffset.left);
			}


			//Draw leader line
			ctx.lineWidth = 2;
			ctx.beginPath();
			ctx.moveTo(popupOffset.left, popupOffset.top);
			ctx.lineTo(xHitPoint, yHitPoint);
			if(drawHorizontal)
				ctx.lineTo(targetOffset.left, targetOffset.top);
			ctx.stroke();
		}
		
	}

	function selectFirstHeader(header) {
		this.detail.target = $('.header > [data-property="' + header[0] + '"]');
		this.detail.launchAtTarget(this.detail.target, 'header');
	}


	return {
		detail : {
			init: init,
			drawCanvasLines: drawCanvasLines,
			launchAtTarget: launchAtTarget,
			selectFirstHeader: selectFirstHeader,

			onDoubleClick: onDoubleClick,
			onDrag: onDrag,
			onModalHide: onModalHide,
			onModalSaveHeader: onModalSaveHeader,
			onModalSaveCell: onModalSaveCell
		}
	}
});