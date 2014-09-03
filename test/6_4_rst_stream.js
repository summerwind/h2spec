var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.4. RST_STREAM', function () {
  it('Sends a RST_STREAM frame with 0x0 stream identifier', function (done) {
    var message = "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createRstStreamFrame(protocol.CODE_CANCEL);
      frame.streamId = 0x0;
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a RST_STREAM frame on a idle stream', function (done) {
    var message = "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createRstStreamFrame(protocol.CODE_CANCEL);
      frame.streamId = 1;
      conn.socket.write(frame.encode());
    });
  });
});