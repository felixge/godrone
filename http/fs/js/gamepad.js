"use strict";

// Gamepad provides a simple wrapper around the HTML5 Gamepad API.
// see http://www.html5rocks.com/en/tutorials/doodles/gamepad/
window.Gamepad = (function() {
  function Gamepad(options) {
    this._onConnect = options.onConnect || function() {};
    this._onClose = options.onClose || function() {};
    this._onChange = options.onChange || function() {};
    this._connected = false;
    this._timestamp = null;
  }

  Gamepad.prototype.connect = function() {
    this._poll();
  };

  Gamepad.prototype._poll = function() {
    var self = this;
    requestAnimationFrame(this._poll.bind(this));

    var gamepad = navigator.webkitGetGamepads()[0];
    if (!gamepad) {
      if (this._connected) {
        this._onClose();
        this._connected = false;
      }
      return;
    }

    if (!this._connected) {
      this._onConnect();
      this._connected = true;
    }

    if (gamepad.timestamp != this._timestamp) {
      this._onChange($.extend(true, {}, gamepad));
    }
    this._timestamp = gamepad.timestamp;
  };

  return Gamepad;
})();
