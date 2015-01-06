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
/* Make it hide itself again if you click on it, so you can see the graphs. */
emdiv.onclick = function() {
  emdiv.style.display = "none";
}

function showEmergency(why) {
  emdiv.style.display = "inherit";
  document.getElementById("emwhy").innerHTML = why
}

function Conn(options) {
  var ws = new WebSocket(options.url);
  var first = true;
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

    // For each response, set our desired alt to the drone's current desired alt,
    // so that we can:
    // - reconnect to a flying drone
    // - monitor a drone that's doing something on it's own (like landing)
    desired.pitch = data.desired.Pitch;
    desired.roll = data.desired.Roll;
    desired.yaw = data.desired.Yaw;
    desired.altitude = data.desired.Altitude;

    if (data.cutout) {
      emergency = true
      showEmergency("Drone cutout: " + data.cutout)
    }

    var latency = Date.now() - lastSend;
    var timeout = Math.max(0, 1000/options.hz-latency);
    lastSend = Date.now();
    setTimeout(function() {
      var msg = {};
      var newDesired = { pitch: desired.pitch, roll: desired.roll,
                         yaw: desired.yaw, altitude: desired.altitude };
      if (isDown[KEYS.c]) {
        msg.calibrate = true;
      }
      if (emergency) {
        newDesired = {pitch: 0, roll: 0, yaw: 0, altitude: 0};
      } else {
        var altSpeed = 0.5/options.hz;
        var yawSpeed = 10/options.hz;
        var speed = 3;
        if (isDown[KEYS.w]) {
          newDesired.pitch = -speed;
        } else if (isDown[KEYS.s]) {
          newDesired.pitch = speed;
        } else {
          newDesired.pitch = 0;
        }
        if (isDown[KEYS.a]) {
          newDesired.roll = -speed;
        } else if (isDown[KEYS.d]) {
          newDesired.roll = speed;
        } else {
          newDesired.roll = 0;
        }
        if (isDown[KEYS.left]) {
          newDesired.yaw += speed;
        } else if (isDown[KEYS.right]) {
          newDesired.yaw -= speed;
        }
        if (isDown[KEYS.up]) {
          newDesired.altitude += altSpeed;
          newDesired.altitude = Math.max(newDesired.altitude, minAltitude);
        } else if (isDown[KEYS.down]) {
          newDesired.altitude -= altSpeed;
          if (newDesired.altitude < minAltitude) {
            newDesired.altitude = 0;
          }
        }
        if (isDown[KEYS.l]) {
            msg.land = true;
        }
      }
      if (emergency || newDesired.pitch != desired.pitch ||
          newDesired.roll != desired.roll ||
          newDesired.yaw != desired.yaw ||
          newDesired.altitude != desired.altitude) {
        msg.setDesired = newDesired;
        desired = newDesired;
      }
      ws.send(JSON.stringify(msg));
    }, timeout);
  };
}

var drone = 'ws://192.168.1.1'
if (window.location.search != "") {
    drone = "ws://"+window.location.search.substring(1)
}

new Conn({
  url: drone,
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
  l: 76,
  p: 80,
    question: 191,
};
var isDown = {};

window.onkeydown = function(e) {
  isDown[e.keyCode] = true;
  switch (e.keyCode) {
    case KEYS.question:
      var helpdiv = document.getElementById("help")
      if (helpdiv.style.display == "none") {
	  helpdiv.style.display = "inherit";
      } else {
	  helpdiv.style.display = "none";
      }
      break;
    case KEYS.esc:
      emergency = true;
      showEmergency("user requested");
      break;
    case KEYS.p:
      pause = !pause;
      break;
  }
}

window.onkeyup = function(e) {
  isDown[e.keyCode] = false;
}
