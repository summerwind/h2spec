var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.5. SETTINGS', function () {
  it('Sends a SETTINGS frame', function (done) {
    var message = 'the endpoint MUST sends a SETTINGS frame with ACK.';

    var conn = helper.createConnection(function () {
      conn.once('frame', function (frame) {
        if (frame.type === protocol.FRAME_TYPE_SETTINGS) {
          assert.ok(frame.ack, message);
          done();
        }
      });

      var frame = helper.createSettingsFrame();
      frame.setMaxConcurrentStreams(100);
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a SETTINGS frame that is not a zero-length with ACK flag', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_FRAME_SIZE_ERROR, message);
        done();
      });

      var frame = helper.createSettingsFrame();
      frame.length = 1;
      frame.ack = true;
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a SETTINGS frame with the stream identifier that is not 0x0', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createSettingsFrame();
      frame.streamId = 0x3;
      frame.setMaxConcurrentStreams(100);
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a incomplete SETTINGS frame', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = new Buffer([
        0x00, 0x00, 0x02, 0x04, 0x04, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x64
      ]);
      conn.socket.write(frame);
    });
  });

  describe('6.5.2.  Defined SETTINGS Parameters', function () {
    describe('SETTINGS_ENABLE_PUSH (0x2)', function () {
      it('Sends the value other than 0 or 1', function (done) {
        var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

        var conn = helper.createConnection(function () {
          conn.once('error', function (err) {
            assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
            done();
          });

          var frame = new Buffer([
            0x00, 0x00, 0x06, 0x04, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x02
          ]);
          conn.socket.write(frame);
        });
      });
    });

    describe('SETTINGS_INITIAL_WINDOW_SIZE (0x4)', function () {
      it('Sends the value above the maximum flow control window size', function (done) {
        var message = 'the endpoint MUST respond with a connection error of type FLOW_CONTROL_ERROR.';

        var conn = helper.createConnection(function () {
          conn.once('error', function (err) {
            assert.equal(err.code, protocol.CODE_FLOW_CONTROL_ERROR, message);
            done();
          });

          var frame = new Buffer([
            0x00, 0x00, 0x06, 0x04, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x04, 0x80, 0x00, 0x00, 0x00
          ]);
          conn.socket.write(frame);
        });
      });
    });

    describe('SETTINGS_MAX_FRAME_SIZE (0x5)', function () {
      it('Sends the value below the initial value', function (done) {
        var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

        var conn = helper.createConnection(function () {
          conn.once('error', function (err) {
            assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
            done();
          });

          var frame = new Buffer([
            0x00, 0x00, 0x06, 0x04, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x05, 0x00, 0x00, 0x3f, 0xff
          ]);
          conn.socket.write(frame);
        });
      });

      it('Sends the value above the maximum allowed frame size', function (done) {
        var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

        var conn = helper.createConnection(function () {
          conn.once('error', function (err) {
            assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
            done();
          });

          var frame = new Buffer([
            0x00, 0x00, 0x06, 0x04, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x05, 0x01, 0x00, 0x00, 0x00
          ]);
          conn.socket.write(frame);
        });
      });
    });
  });
});