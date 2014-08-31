var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.7. PING', function () {
  it('Sends a PING frame', function (done) {
    var message = 'the endpoint MUST sends a PING frame with ACK.';

    var conn = helper.createConnection(function () {
      conn.once('frame', function (frame) {
        if (frame.type === protocol.FRAME_TYPE_PING) {
          assert.ok(frame.ack, message);
          done();
        }
      });

      var frame = helper.createPingFrame();
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a PING frame with the stream identifier that is not 0x0', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createPingFrame();
      frame.streamId = 0x3;
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a PING frame with a length field value other than 8', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_FRAME_SIZE_ERROR, message);
        done();
      });

      var frame = helper.createPingFrame();
      frame.length = 1;
      conn.socket.write(frame.encode());
    });
  });
});