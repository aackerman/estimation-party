define(
  [
    "jquery"
  ],

  function(
    $
  ) {
    var room = window.location.pathname.replace('\/room\/', '')
    var ws = new WebSocket("ws://localhost:8000/ws/" + room);
    ws.onclose = function(e) {
      console.log("WebSocket closed", e);
    };
    ws.onerror = function(e) {
      console.log("WebSocket error", e);
    };
    ws.onmessage = function(e) {
      var data = JSON.parse(e.data);
      socket.route(data.route, data.data);
    };

    var pubsub = $({});

    var socket = {
      __callbacks: {},
      emit: function(route, data) {
        ws.send(JSON.stringify({route: route, data: data}));
      },
      on: function(r, fn) {
        socket.__callbacks[r] = fn
      },
      route: function() {
        var args = [].slice.call(arguments);
        var route = args.shift()
        if (socket.__callbacks[route]) {
          socket.__callbacks[route].apply(null, args);
        } else {
          console.log("ROUTE DOES NOT EXIST", route)
        }
      }
    };

    return socket;
  }
);
