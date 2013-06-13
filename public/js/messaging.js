define(
  [
    "jquery"
  ],

  function($) {
    var Messaging = {
      voting: function(bool) {
        var s = (bool ? 'open' : 'closed')
        $('.voting').text('Voting is ' + s);
      },
      notification: function(msg) {
        $('.notifications').text(msg);
      },
      ticket: function(t) {
        console.log('set ticket', t);
        $('.ticket').text(t);
      }
    };

    return Messaging;
  }

)
