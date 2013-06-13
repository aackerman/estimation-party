define(
  [
    "jquery"
  ],

  function(
    $
  ) {

    var ws = new WebSocket("ws://localhost:8000/ws");
    ws.onclose = function(e) {
      console.log(e);
    };
    ws.onerror = function(e) {
      console.log(e);
    };
    ws.onmessage = function(e) {
      var data = JSON.parse(e.data);
      console.log(e)
      socket.emit(data.route);
    };

    var pubsub = $({});

    var socket = {
      send: function(route, data) {
        var json = JSON.stringify({route: route, data: data})
        console.log('json sent', json);
        ws.send(json);
      },
      on: function() {
        pubsub.on.apply(pubsub, arguments);
      },
      off: function() {
        pubsub.off.apply(pubsub, arguments);
      },
      emit: function() {
        pubsub.trigger.apply(pubsub, arguments);
      }
    };

    return socket;
  }
);
