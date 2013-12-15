"use strict";

// Gamepad provides a simple wrapper around the HTML5 Gamepad API.
// see http://www.html5rocks.com/en/tutorials/doodles/gamepad/
window.Gamepad = (function() {
  function Gamepad(options) {
    this._onConnect = options.onConnect;
    this._onClose = options.onClose;
    this._onChange = options.onChange;
    this._connected = false;
    this._timestamp = null;
  }

  Gamepad.prototype.connect = function() {
    this._poll();
  };

  Gamepad.prototype._poll = function() {
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
      this._onChange(gamepad);
    }
    this._timestamp = gamepad.timestamp;
  };

  return Gamepad;
})();
