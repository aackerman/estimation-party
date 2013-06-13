define(
  [
    "jquery",
    "socket",
    "messaging"
  ],

  function(
    $,
    socket
  ){
    $('.estimate').on('click', function(e){
      console.log('click')
      var $el = $(e.target);
      socket.send('vote', { points: $el.data('value').toFixed(1).toString() });
      $('.estimate').off('click');
    });

    $('.start-voting').on('click', function() {
      socket.send('start', { ticket: $('#ticket').val() })
    });

    $('.controls').show();

    socket.on('start', function(data){
      messaging.ticket(data.ticket)
    });
  }
);
