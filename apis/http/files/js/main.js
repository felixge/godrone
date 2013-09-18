$(function() {
  var uri = 'ws://'+window.location.host+'/ws';
  console.log('Connecting websocket: '+uri);
  var socket = new WebSocket(uri);
  socket.onopen = function() {
    console.log('Websocket connection established');
  };
  socket.onmessage = function(event) {
    console.log('New websocket message', event.data);
  };
});
