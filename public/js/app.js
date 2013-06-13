define(
  [
    "jquery",
    "socket"
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
  }
);
