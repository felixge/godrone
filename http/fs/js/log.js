"use strict";

// Log implements a simple logger. This can be extended with more functionality
// e.g. logging into a DOM element in the future.
window.Log = (function() {
  var levels = [
    'debug',
    'info',
    'warn',
    'error',
    'fatal',
  ];

  function Log(level) {
    this._num = levelNum(level);
  }

  // levelNum returns the numeric level for the given level string
  function levelNum(level) {
    var i = levels.indexOf(level);
    if (i < 0) {
      throw new Error('invalid level: '+level);
    }
    return i;
  }

  // logFn returns a function that handles log calls for the given level
  function logFn(level) {
    return function() {
      var num = levelNum(level);
      if (num < this._num) {
        // we're not supposed to log at this level
        return;
      }

      var args = Array.prototype.slice.call(arguments);
      args[0] = '['+level+'] ' +args[0];
      console.log.apply(console, args);
    };
  }

  // add methods for each log level
  levels.forEach(function(level) {
    Log.prototype[level] = logFn(level);
  });

  return Log;
})()
