define(
  [
    "jquery"
  ],

  function(
    $
  ){
    function route(e) {
      console.log(e);
    }

    var ws = new WebSocket("ws://localhost:8000/ws");
    ws.onclose = route;
    ws.onerror = route;
    ws.onmessage = route;

    var pubsub = $({});

    var socket = {
      send: function(str) {
        ws.send(str);
      },
      on: function() {
        pubsub.on.apply(pubsub, arguments);
      },
      off: function() {
        pubsub.off.apply(pubsub, arguments);
      }
    };

    ws.onopen = function() {
      console.log('onopen')
      ws.send(JSON.stringify({route: 'hello', msg: 'world'}));
    }

    return socket;
  }
);
