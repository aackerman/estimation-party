define(
  [
    "jquery",
    "lodash",
    "messaging",
    "socket",
    "estimate",
    "timing",
    "controls",
    "state"
  ],

  function(
    $,
    _,
    messaging,
    socket,
    estimate,
    timing,
    controls,
    state
  ) {
    var App = {
      results: function(data) {
        timing.end();
        messaging.voting(false);
        estimate.results(data);
      },
      domEvents: function() {
        $('.start-voting').on('click', function() {
          socket.emit('start', { ticket: $('#ticket').val() })
        });
      },
      start: function(data) {
        state.reset();
        timing.start(data.timing);
        messaging.voting(true);
        messaging.ticket(data.ticket);
        messaging.notification('Start Voting!');
        $('.estimate')
          .on('click', estimate.vote)
          .removeClass('voted');
      },
      socketEvents: function() {
        socket.on('reset', App.reset);
        socket.on('sync', App.sync)
        socket.on('results', App.results);
        socket.on('start', App.start);
      },
      sync: function(data) {
        console.log('sync data', data);
        messaging.ticket(data.ticket);
        messaging.voting(data.voting);
        if (data.voting) App.start(data);
      },
      init: function() {
        state.reset();
        controls.init();
        App.domEvents();
        App.socketEvents();
        return App;
      }
    };
    return App.init();
  }
);
