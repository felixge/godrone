"use strict";

var connection = new Connection({url: "ws://192.168.1.1/ws"});
var gamepad = new Gamepad();

(function main() {
  $.when(connection.start(), gamepad.start()).then(function() {
    gamepad.onchange = function(state) {
      var throttle = Math.max(-filter(state.axes[3]), 0);
      if (state.buttons[0]) {
        throttle = 0.5;
      }

      var control = {
        Pitch: filter(state.axes[1])*6,
        Roll: filter(state.axes[0])*6,
        Yaw: filter(state.axes[2])*6,
        Throttle: throttle,
      }
      console.log(state.axes, control);
      connection.send(control);
    };
  }, function() {
    console.log('error', arguments);
  });
})();

function filter(val) {
  return (Math.abs(val) < 0.05)
    ? 0
    : val
}

function pollGamepad() {
  var gamepad = navigator.webkitGetGamepads()[0];
  if (!gamepad) {
    requestAnimationFrame(pollGamepad);
    return;
  }

  requestAnimationFrame(pollGamepad);
}
