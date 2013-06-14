define(
  [
    "jquery"
  ],

  function($){

    return {
      init: function() {
        if(window.location.hash == '#master') this.show();
      },
      show: function(){
        $('.controls').show();
      }
    }
  }
)
