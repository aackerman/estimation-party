define(
  [
    "jquery",
    "messaging"
  ],

  function(
    $,
    messaging
  ) {

    var timing = {
      total: 300,
      start: function(time) {
        var $el = $('.progress .bar');
        $el.css({ width: "100%" });
        $el
          .addClass('active')
          .animate({ width: 0 }, timing.total * 1e3, 'linear');
      },
      end: function() {
        $('.progress .bar').css({ width: 0 }).stop();
      }
    };

    return timing;
  }

)
