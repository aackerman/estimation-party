define(
  [
    "jquery"
  ],

  function(
    $
  ) {
    function route(e) {
      console.log(e);
    }

    var ws = new WebSocket("ws://localhost:8000/ws");
    ws.onclose = route;
    ws.onerror = route;
    ws.onmessage = route;

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
      }
    };

    return socket;
  }
);
