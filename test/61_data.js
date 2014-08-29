var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.1. DATA', function () {
  it('Sends a DATA frame with 0x0 stream identifier', function (done) {
    var message = "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var frame = helper.createDataFrame('test', true);
      frame.streamId = 0x0;
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a DATA frame on the stream that is not opend', function (done) {
    var message = "The endpoint MUST respond with a stream error of type STREAM_CLOSED.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_STREAM_CLOSED, message);
        done();
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

      var dataFrame = helper.createDataFrame('test', true);
      conn.socket.write(dataFrame.encode());
    });
  });

  it('Sends a DATA frame with invalid pad length', function (done) {
    var message = "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var headers = [
        [ ':method',    'POST'           ],
        [ ':scheme',    'http'           ],
        [ ':path',      '/'              ],
        [ ':authority', helper.getHost() ]
      ];

      var context = helper.createCompressionContext();
      var fragment = context.compress(headers);

      var headersFrame = helper.createHeadersFrame(fragment, true, false);
      conn.socket.write(headersFrame.encode());

      var dataFrame = helper.createDataFrame('test', true);
      dataFrame.setPadding(10);
      dataFrame.length -= 10;
      conn.socket.write(dataFrame.encode());
    });
  });
});