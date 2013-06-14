define(
  [
    "jquery",
    "lodash",
    "socket",
    "messaging",
    "state"
  ],

  function($, _, s, messaging, state) {
    function winner(data) {
      var max = -1, top;
      _.forEach(data, function(v, k){if (v > max) top = k});
      return top
    }

    return {
      // called by a DOM event
      vote: function(e) {
        var $el = $(e.target), points;
        if (!state.get('voted')) {
          $(e.target).addClass('voted');
          $('.estimate').off('click');
          state.set('voted', true);
          points = $el.data('value').toString();
          s.emit('vote', { points: points });
          messaging.notification('Voted for ' + points);
        }
      },
      results: function(data){
        var points, height;
        var top = winner(data);
        console.log(top)
        $('.estimate').removeClass('voted')
        $('.estimate').each(function(i, el) {
          $el = $(el);
          points = $el.data('value');
          height = data[points] * 10 + 30;
          console.log(points, top, points == top)
          if (points == top) $el.addClass('winner');
          $el.height(height + 'px');
        });
      }
    };
  }

)
