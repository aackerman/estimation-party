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
      return _.max(_.values(data));
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
        $('.estimate').removeClass('voted')
        $('.estimate').each(function(i, el) {
          $el = $(el);
          points = $el.data('value');
          height = data[points] * 10 + 30;
          if (points == winner(data)) $el.addClass('winner');
          $el.height(height + 'px');
        });
      }
    };
  }

)
