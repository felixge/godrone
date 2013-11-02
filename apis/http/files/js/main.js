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
  };

  var lastHz = Date.now();
  setInterval(function() {
    console.log('navdata rate: %s hz', count/((lastHz-Date.now())/1000), navdata, fields);
    count = 0;
    lastHz = Date.now();
  }, 1000);

  //
  // We use an inline data source in the example, usually data would
  // be fetched from a server

  var totalPoints = 300, scale = 100;
  var fields = {};

  function getNavdata(fieldName) {
    var val = navdata[fieldName];

    var field = fields[fieldName];
    if (!field) {
      var defaultMin, defaultMax;
      if (/^A/.test(fieldName)) {
        defaultMin = 0;
        defaultMax = 4096;
      } else if (/^G/.test(fieldName)) {
        defaultMin = -10*1000;
        defaultMax = 10*1000;
      }

      field = fields[fieldName] = {min: defaultMin, max: defaultMax, data: []}
    }

    if (field.min === undefined || val < field.min) {
      field.min = val;
    }
    if (field.max === undefined || val > field.max) {
      field.max = val;
    }

    if (field.data.length == 0) {
      for (var i = 0; i < totalPoints; i++) {
        field.data.push(0);
      }
    }

    scaled = ((val - field.min) / (field.max - field.min)) * scale;

    field.data.shift();
    field.data.push(scaled);

    // Zip the generated y values with the x values

    var res = [];
    for (var i = 0; i < field.data.length; ++i) {
      res.push([i, field.data[i]])
    }

    return {
      data: res,
      label: fieldName,
    };
  }

  function getPlots() {
    return [
      getNavdata('Gx'),
      getNavdata('Gy'),
      getNavdata('Gz'),
      getNavdata('Ax'),
      getNavdata('Ay'),
      getNavdata('Az'),
    ];
  }

  var plot;

  function update() {
    if (navdata) {
      if (!plot) {
        plot = $.plot("#plot", getPlots(), {
          series: {
            shadowSize: 0,
          },
          yaxis: {
            min: 0,
            max: scale
          },
          xaxis: {
            show: false
          }
        });
      }

      plot.setData(getPlots());
      plot.draw();
    }
    requestAnimationFrame(update);
  }

  update();
});
