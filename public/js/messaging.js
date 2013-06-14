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
        if (/[0-9]{4}/.test(t)) {
          $('.ticket').html("<a href='http://jira/browse/OF-" + t + "' target='_blank'>" + t + "</a>");
        } else  {
          $('.ticket').text(t);
        }
      }
    };

    return Messaging;
  }

)
