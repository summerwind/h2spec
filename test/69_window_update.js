var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.9. WINDOW_UPDATE', function () {
  it('Sends a WINDOW_UPDATE frame with an flow control window increment of 0', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createWindowUpdateFrame(0);
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a WINDOW_UPDATE frame with an flow control window increment of 0 on a stream', function (done) {
    var message = 'the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.on('frame', function (frame) {
        if (frame.type === protocol.FRAME_TYPE_RST_STREAM) {
          assert.equal(frame.errorCode, protocol.CODE_PROTOCOL_ERROR, message);
          done();
        }
      });

      var headers = [
        [ ':method',    'GET'            ],
        [ ':scheme',    'http'           ],
        [ ':path',      '/'              ],
        [ ':authority', helper.getHost() ]
      ];

      var context = helper.createCompressionContext();
      var fragment = context.compress(headers);

      var headersFrame = helper.createHeadersFrame(fragment, true, true);
      conn.socket.write(headersFrame.encode());

      var windowUpdateframe = helper.createWindowUpdateFrame(0);
      windowUpdateframe.streamId = headersFrame.streamId;
      conn.socket.write(windowUpdateframe.encode());
    });
  });

  describe('6.9.2. Initial Flow Control Window Size', function () {
    it('Sends a SETTINGS_INITIAL_WINDOW_SIZE settings with an exceeded maximum window size value', function (done) {
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
});