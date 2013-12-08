"use strict";

function Gamepad() {
  this._started = null;
  this._timestamp = null;
  this.onchange = function() {};
}

Gamepad.prototype.start = function() {
  console.log('Starting gamepad');
  this._poll();
  this._started = new $.Deferred();
  return this._started.promise();
};

Gamepad.prototype._poll = function() {
  var state = navigator.webkitGetGamepads()[0];
  if (!state) {
    requestAnimationFrame(this._poll.bind(this));
    return;
  }

  if (this._started) {
    console.log('Gamepad started');
    this._started.resolve();
    this._started = null;
  }

  if (state.timestamp !== this._timestamp) {
    this.onchange(state);
    this._timestamp = state.timestamp;
  }
  requestAnimationFrame(this._poll.bind(this));
};
