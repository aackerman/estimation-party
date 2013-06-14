define(
  [
    "jquery",
    "lodash",
    "socket",
    "messaging",
    "state"
  ],

  function($, _, s, messaging, state) {
    function winner(results) {
      var maxvotes = 0, points;
      _.each(results, function(votes, value) {
        if (votes > maxvotes) points = value;
      });
      return points;
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
        var max = _.max(data, function(r){ return r.points });
        console.log(max)

        $('.estimate').each(function(i, el) {
          $el = $(el);
          points = $el.data('value');
          height = data[points] * 20;
          if (points == winner(data)) $el.addClass('winner');
          if (height) $el.addClass('voted');
          $el.height(height + 50 + 'px');
        });
      }
    };
  }

)
