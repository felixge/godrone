$(function() {
  var uri = 'ws://'+window.location.host+'/ws';
  console.log('Connecting websocket: '+uri);
  var socket = new WebSocket(uri);
  socket.onopen = function() {
    console.log('Websocket connection established');
  };

  var count = 0;
  var navdata;
  socket.onmessage = function(e) {
    count++;
    navdata = JSON.parse(e.data);
    navdata.Received = Date.now();
  };

  var lastHz = Date.now();
  setInterval(function() {
    console.log('navdata rate: %s hz', count/((lastHz-Date.now())/1000), navdata, min, max);
    count = 0;
    lastHz = Date.now();
  }, 1000);

  var smoothie = new SmoothieChart({millisPerPixel: 1,minValue:500, maxValue:4000});
  smoothie.streamTo(document.getElementById('sensors'));

  var sensors = ['Az'];
  var lines = {};
  sensors.forEach(function(sensor) {
    var line = new TimeSeries();
    smoothie.addTimeSeries(line);
    lines[sensor] = line;
  });

  var min = undefined;
  var max = undefined;

  function update() {
    if (navdata) {
      sensors.forEach(function(sensor) {
        var val = navdata[sensor];
        if (min === undefined || val < min) {
          min = val;
        }
        if (max === undefined || val > max) {
          max = val;
        }
        lines[sensor].append(navdata.Received, val)
      });
    }

    requestAnimationFrame(update);
  }

  update();
});
