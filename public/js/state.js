define(
  [
    "lodash"
  ],

  function(_) {
    return {
      __default: {
        voting: false
      },
      __attrs: {},
      get: function(attr) {
        return this.__attrs[attr];
      },
      set: function(attr, val) {
        return this.__attrs[attr] = val;
      },
      reset: function(){
        this.__attrs = _.cloneDeep(this.__default);
        $('.estimate').height('30px');
      }
    };
  }
)
