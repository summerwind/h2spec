var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.2. HEADERS', function () {
  it('Sends a HEADERS frame followed by any frame other than CONTINUATION', function (done) {
    var message = "The endpoint MUST treat the receipt of any other type of frame as a connection error of type PROTOCOL_ERROR.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
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

      var headersFrame = helper.createHeadersFrame(fragment, false, false);
      conn.socket.write(headersFrame.encode());

      var dataFrame = helper.createDataFrame('test', true);
      conn.socket.write(dataFrame.encode());
    });
  });

  it('Sends a HEADERS frame followed by a frame on a different stream', function (done) {
    var message = "The endpoint MUST treat the receipt of a frame on a different stream as a connection error of type PROTOCOL_ERROR.";

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
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

      var headersFrame = helper.createHeadersFrame(fragment, false, false);
      conn.socket.write(headersFrame.encode());

      var dataFrame = helper.createDataFrame('test', true);
      dataFrame.streamId = 3;
      conn.socket.write(dataFrame.encode());
    });
  });

  it('Sends a HEADERS frame with 0x0 stream identifier', function (done) {
    var message = "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.";

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

      var frame = helper.createHeadersFrame(fragment, true, false);
      frame.streamId = 0x0;
      conn.socket.write(frame.encode());
    });
  });

  it('Sends a HEADERS frame with invalid pad length', function (done) {
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

      var frame = helper.createHeadersFrame(fragment, true, false);
      frame.setPadding(10);
      frame.length -= 10;
      conn.socket.write(frame.encode());
    });
  });
});