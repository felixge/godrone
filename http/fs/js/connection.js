"use strict";

function Connection(options) {
  this._url = options.url;
  this._conn = null;
}

Connection.prototype.start = function() {
  if (this._conn) {
    return;
  }

  console.log('Starting websocket connection to '+this._url);

  var d = new $.Deferred();
  this._conn = new WebSocket(this._url);
  this._conn.onopen = function(e) {
    console.log('Websocket connected');
    d.resolve();
  };
  this._conn.onerror = function(e) {
    d.reject(e);
  }
  return d.promise();
};

Connection.prototype.send = function(data) {
  this._conn.send(JSON.stringify(data));
};
