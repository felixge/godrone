$(function() {
  var conn = new WebSocket('ws://10.0.1.25/ws');
  conn.binaryType = 'arraybuffer';
  conn.onopen = function () {
    console.log('opened');
  };
  conn.onerror = function (error) {
    console.log('WebSocket Error ' + error);
  };
  conn.onclose = function () {
    console.log('close', arguments);
  };

  var divs = [];

  conn.onmessage = function (e) {
    console.log(e.data);
  };

  var $speed = $('#js_speed');
  $speed.change(function() {
    conn.send($speed.val())
  });
});
