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


  var $canvas =  $('canvas');
  var width = $canvas.width();
  var height = $canvas.height();
  var ctx = window.ctx = $canvas[0].getContext('2d');
  ctx.save();

  var points = [];
  var maxPoints = 500;

  var colors = {Roll: 'red', Pitch: 'blue'};

  var j = 0;
  conn.onmessage = function (e) {
    j++;
    var data = JSON.parse(e.data);
    if (j % 100 === 0) {
      console.log(data.Roll, data.Pitch);
    }

    points.push(data);

    if (points.length > maxPoints) {
      points.shift();
    }

    ctx.clearRect(0, 0, width, height);
    for (var i = 0; i < points.length; i++) {
      var x = (width / maxPoints) * i;

      ['Roll', 'Pitch'].forEach(function(field) {
        var y = (points[i][field] / 180) * height + 180;
        ctx.beginPath();
        ctx.arc(x, y, 1, 0, 2 * Math.PI, false);
        ctx.fillStyle = colors[field];
        ctx.fill();
      });
    }
    //console.log(data);
  };

  var $speed = $('#js_speed');
  $speed.change(function() {
    conn.send($speed.val())
  });
});
