var net = require('net'),
    assert = require('assert');

var helper = require('./helper');

describe('3.5. HTTP/2 Connection Preface', function () {
  it('Sends invalid connection preface', function (done) {
    var message = 'The endpoint MUST terminate the TCP connection.';
    var send = false;

    var options = {
      port: helper.getPort(),
      host: helper.getHost()
    };

    var socket = net.connect(options, function () {
      socket.write('INVALID CONNECTION PREFACE');
      send = true;
    });

    socket.once('close', function () {
      assert.ok(send, message);
      done();
    });
  });
});