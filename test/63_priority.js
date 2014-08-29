var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.3. PRIORITY', function () {
  it('Sends a PRIORITY frame with 0x0 stream identifier', function (done) {
    var message = "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createPriorityFrame(1, 10, true);
      frame.streamId = 0x0;
      conn.socket.write(frame.encode());
    });
  });
});