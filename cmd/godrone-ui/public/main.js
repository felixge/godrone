"use strict";

var minAltitude = 0.3;
var pause = false;
var emergency = false;
var colors = ['red', 'black'];
var desired = {pitch: 0, roll: 0, yaw: 0, altitude: 0};
var charts = [
  {
    title: 'Altitude',
    labels: ['Time', 'Altitude', 'Desired'],
    data: [],
    colors: colors,
  },
  {
    title: 'Pitch',
    labels: ['Time', 'Actual', 'Desired'],
    data: [],
    valueRange: [-90, 90],
    colors: colors,
  },
  {
    title: 'Roll',
    labels: ['Time', 'Actual', 'Desired'],
    data: [],
    valueRange: [-90, 90],
    colors: colors,
  },
  {
    title: 'Yaw',
    labels: ['Time', 'Actual', 'Desired'],
    data: [],
    valueRange: [-90, 90],
    colors: colors,
  },
];

var chartData = [];
var graphs = [];
for (var i = 0; i < charts.length; i++) {
  graphs[i] = new Dygraph(document.getElementById('chart'+i), charts[i].data, {
  	drawPoints: true,
  	showRoller: true,
  	labels: charts[i].labels,
  	valueRange: charts[i].valueRange,
  	title: charts[i].title,
  	colors: charts[i].colors,
	})
}

var emdiv = document.getElementById("emergency");

function Conn(options) {
  var ws = new WebSocket(options.url);
  var reconnect = function() {
    ws.onopen = function() {};
    ws.onmessage = function() {};
    ws.onerror = function() {};
    reconnect = function() {};
    ws.close();
    setTimeout(function() {
      new Conn(options);
    }, 1000);
  };
  var lastSend;

  ws.onerror = function(e) {
    console.log('Failed to connect to: '+options.url, e);
    reconnect();
  };
  ws.onopen = function() {
    console.log('Connected to: '+options.url);
    ws.send(JSON.stringify({}));
    lastSend = Date.now();
  };
  ws.onmessage = function(e) {
    try {
      var data = JSON.parse(e.data);
    } catch (err) {
      console.log('Received bad ws data', err);
      reconnect();
      return;
    }
    var time = new Date(data.time);
    charts[0].data.push([
      time,
      data.actual.Altitude,
      data.desired.Altitude,
    ]);
    charts[1].data.push([
      time,
      data.actual.Pitch,
      data.desired.Pitch,
    ]);
    charts[2].data.push([
      time,
      data.actual.Roll,
      data.desired.Roll,
    ]);
    charts[3].data.push([
      time,
      data.actual.Yaw,
      data.desired.Yaw,
    ]);
    if (!pause) {
	for (var i = 0; i < graphs.length; i++) {
          graphs[i].updateOptions({file: charts[i].data, dateWindow: [time-10000, time]});
	}
    }

    var latency = Date.now() - lastSend;
    var timeout = Math.max(0, 1000/options.hz-latency);
    lastSend = Date.now();
    setTimeout(function() {
      var msg = {setDesired: desired};
      if (emergency) {
        desired = {pitch: 0, roll: 0, yaw: 0, altitude: 0};
      } else {
        var altSpeed = 0.5/options.hz;
        var yawSpeed = 10/options.hz;
        var speed = 3;
        if (isDown[KEYS.w]) {
          desired.pitch = -speed;
        } else if (isDown[KEYS.s]) {
          desired.pitch = speed;
        } else {
          desired.pitch = 0;
        }
        if (isDown[KEYS.a]) {
          desired.roll = -speed;
        } else if (isDown[KEYS.d]) {
          desired.roll = speed;
        } else {
          desired.roll = 0;
        }
        if (isDown[KEYS.left]) {
          desired.yaw += speed;
        } else if (isDown[KEYS.right]) {
          desired.yaw -= speed;
        }
        if (isDown[KEYS.up]) {
          desired.altitude += altSpeed;
          desired.altitude = Math.max(desired.altitude, minAltitude);
        } else if (isDown[KEYS.down]) {
          desired.altitude -= altSpeed;
          if (desired.altitude < minAltitude) {
            desired.altitude = 0;
          }
        }
        if (isDown[KEYS.c]) {
          msg.calibrate = true;
        }
      }
      ws.send(JSON.stringify(msg));
    }, timeout);
  };
}

new Conn({
  url: 'ws://192.168.1.1',
  hz: 30,
});

var KEYS = {
  esc: 27,
  left: 37,
  up: 38,
  right: 39,
  down: 40,
  c: 67,
  w: 87,
  s: 83,
  a: 65,
  d: 68,
  p: 80,
};
var isDown = {};

window.onkeydown = function(e) {
  console.log(e.keyCode);
  isDown[e.keyCode] = true;
  switch (e.keyCode) {
    case KEYS.esc:
      emergency = true;
      emdiv.style.display = "inherit";
      break;
    case KEYS.p:
      pause = !pause;
      break;
  }
}

window.onkeyup = function(e) {
  isDown[e.keyCode] = false;
}
