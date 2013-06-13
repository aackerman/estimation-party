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
      total: 20,
      start: function(time) {
        var $el = $('.progress .bar');
        $el.css({ width: "100%" });
        var t = (time || timing.total) * 1e3
        $el
          .addClass('active')
          .animate({ width: "0" }, t, 'linear');
      },
      end: function() {
        $('.progress .bar').css({ width: 0 }).stop();
      }
    };

    return timing;
  }

)
