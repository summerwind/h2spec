var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.8. GOAWAY', function () {
  it('Sends a GOAWAY frame with the stream identifier that is not 0x0', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createGoawayFrame(1, protocol.CODE_NO_ERROR, 'h2check');
      frame.streamId = 0x1;
      conn.socket.write(frame.encode());
    });
  });
});