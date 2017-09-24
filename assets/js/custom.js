$(document).ready(function() {

	// Variables
	var $codeSnippets = $('.code-example-body'),
		$nav = $('.navbar'),
		$body = $('body'),
		$window = $(window),
		$popoverLink = $('[data-popover]'),
		navOffsetTop = $nav.offset().top,
		$document = $(document),
		entityMap = {
			"&": "&amp;",
			"<": "&lt;",
			">": "&gt;",
			'"': '&quot;',
			"'": '&#39;',
			"/": '&#x2F;'
		}

	function init() {
		$window.on('scroll', onScroll)
		$window.on('resize', resize)
		$('a[href^="#"]').on('click', smoothScroll)
		buildSnippets();
	}

	function smoothScroll(e) {
		e.preventDefault();
		$(document).off("scroll");
		var target = this.hash,
			menu = target;
		$target = $(target);
		$('html, body').stop().animate({
			'scrollTop': $target.offset().top-40
		}, 0, 'swing', function () {
			window.location.hash = target;
			$(document).on("scroll", onScroll);
		});
	}

	$("#button").click(function() {
		$('html, body').animate({
			scrollTop: $("#elementtoScrollToID").offset().top
		}, 2000);
	});

	function resize() {
		$body.removeClass('has-docked-nav')
		navOffsetTop = $nav.offset().top
		onScroll()
	}

	function onScroll() {
		if(navOffsetTop < $window.scrollTop() && !$body.hasClass('has-docked-nav')) {
			$body.addClass('has-docked-nav')
		}
		if(navOffsetTop > $window.scrollTop() && $body.hasClass('has-docked-nav')) {
			$body.removeClass('has-docked-nav')
		}
	}

	function buildSnippets() {
		$codeSnippets.each(function() {
			var newContent = escapeHtml($(this).html())
			$(this).html(newContent)
		})
	}

	var newIn = '<div id="add-point">' +
	'<div class="form-group">' + 
	'<label>Point </label><small> Add debate points and counterpoints. Each response should be one cohesive thought. If possible it should be backed up by data and links etc...</small>' +
	'<textarea class=" form-control" id="point-Topic" name="Topic" rows="3"></textarea>' + 
	'</div>' +
	'<button class="btn btn-success" role="submit">Add</button>' + 
	'</div>';
	$("#add-input").click(function(e){
		if (newIn) {
			$('#add-input').after(newIn);
			$('#add-input').html('<span class="glyphicon glyphicon-remove" aria-hidden="true"></span>');
			newIn = null;
		}else{
			newIn = $('#add-point').detach(); 
			$('#add-input').html('<span class="glyphicon glyphicon-plus" aria-hidden="true"> Reply</span>');
		}
	});

	init();

});
