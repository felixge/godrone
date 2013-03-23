var http = require('http');
var send = require('send');
var url = require('url');
var port = process.env.PORT || 1080;

http.createServer(function(req, res){
  send(req, url.parse(req.url).pathname)
    .root(__dirname + '/public')
    .pipe(res);
}).listen(port, function() {
  console.log('debug-client listening on http port %d', port);
});
