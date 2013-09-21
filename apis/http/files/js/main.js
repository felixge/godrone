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
    console.log('navdata rate: %s hz', count/((lastHz-Date.now())/1000), navdata);
    count = 0;
    lastHz = Date.now();
  }, 1000);

  //
		// We use an inline data source in the example, usually data would
		// be fetched from a server

		var totalPoints = 300, scale = 100;
    var fields = {};

		function getNavdata(fieldName) {
      var val = (navdata)
        ? navdata[fieldName]
        : 0;

      var field = fields[fieldName];
      if (!field) {
        field = fields[fieldName] = {min: 0, max: 0, data: []}
      }

      if (val < field.min) {
        field.min = val;
      }
      if (val > field.max) {
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

			return res;
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

		var plot = $.plot("#placeholder", getPlots(), {
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

		function update() {
			plot.setData(getPlots());

			// Since the axes don't change, we don't need to call plot.setupGrid()

			plot.draw();
			requestAnimationFrame(update);
		}

		update();
});
