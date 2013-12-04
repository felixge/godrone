"use strict";

var connection = new Connection({url: "ws://192.168.1.1/ws"});
var gamepad = new Gamepad();

(function main() {
  $.when(connection.start(), gamepad.start()).then(function() {
    gamepad.onchange = function(state) {
      console.log(state.axes);
      connection.send({
        Pitch: state.axes[1],
        Roll: state.axes[0],
        Yaw: state.axes[2],
        Vertical: state.axes[3],
      });
    };
  }, function() {
    console.log('error', arguments);
  });
})();


function pollGamepad() {
  var gamepad = navigator.webkitGetGamepads()[0];
  if (!gamepad) {
    requestAnimationFrame(pollGamepad);
    return;
  }

  requestAnimationFrame(pollGamepad);
}
