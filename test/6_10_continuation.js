var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('6.10. CONTINUATION', function () {
  it('Sends a CONTINUATION frame', function (done) {
    var message = 'The endpoint must accept the frame.';

    var conn = helper.createConnection(function () {
      conn.once('frame', function (frame) {
        assert.equal(frame.type, protocol.FRAME_TYPE_HEADERS, message);
        done();
      });

      var headers = [
        [ ':method',    'GET'                      ],
        [ ':scheme',    'http'                     ],
        [ ':path',      '/'                        ],
        [ ':authority', helper.getHost()           ],
        [ 'x-dummy1',   helper.getDummyData(10000) ],
        [ 'x-dummy2',   helper.getDummyData(10000) ]
      ];

      var context = helper.createCompressionContext();
      var fragment = context.compress(headers);

      var headersFrame = helper.createHeadersFrame(fragment.slice(0, 16384), false, true);
      conn.socket.write(headersFrame.encode());

      var continuationFrame = helper.createContinuationFrame(fragment.slice(16384), true);
      conn.socket.write(continuationFrame.encode());
    });
  });

  it('Sends multiple CONTINUATION frames', function (done) {
    var message = 'The endpoint must accept the frames.';

    var conn = helper.createConnection(function () {
      conn.once('frame', function (frame) {
        assert.equal(frame.type, protocol.FRAME_TYPE_HEADERS, message);
        done();
      });

      var headers = [
        [ ':method',    'GET'                      ],
        [ ':scheme',    'http'                     ],
        [ ':path',      '/'                        ],
        [ ':authority', helper.getHost()           ],
        [ 'x-dummy1',   helper.getDummyData(10000) ],
        [ 'x-dummy2',   helper.getDummyData(10000) ],
        [ 'x-dummy3',   helper.getDummyData(10000) ],
        [ 'x-dummy4',   helper.getDummyData(10000) ],
        [ 'x-dummy5',   helper.getDummyData(10000) ]
      ];

      var context = helper.createCompressionContext();
      var fragment = context.compress(headers);

      var headersFrame = helper.createHeadersFrame(fragment.slice(0, 16384), false, true);
      conn.socket.write(headersFrame.encode());

      var continuationFrame1 = helper.createContinuationFrame(fragment.slice(16384, 32767), false);
      conn.socket.write(continuationFrame1.encode());

      var continuationFrame2 = helper.createContinuationFrame(fragment.slice(32767), true);
      conn.socket.write(continuationFrame2.encode());
    });
  });

  it('Sends a CONTINUATION frame followed by any frame other than CONTINUATION', function (done) {
    var message = 'The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var headers = [
        [ ':method',    'GET'                      ],
        [ ':scheme',    'http'                     ],
        [ ':path',      '/'                        ],
        [ ':authority', helper.getHost()           ],
        [ 'x-dummy1',   helper.getDummyData(10000) ],
        [ 'x-dummy2',   helper.getDummyData(10000) ],
        [ 'x-dummy3',   helper.getDummyData(10000) ],
        [ 'x-dummy4',   helper.getDummyData(10000) ],
        [ 'x-dummy5',   helper.getDummyData(10000) ]
      ];

      var context = helper.createCompressionContext();
      var fragment = context.compress(headers);

      var headersFrame = helper.createHeadersFrame(fragment.slice(0, 16384), false, true);
      conn.socket.write(headersFrame.encode());

      var continuationFrame = helper.createContinuationFrame(fragment.slice(16384, 32767), false);
      conn.socket.write(continuationFrame.encode());

      var dataFrame = helper.createDataFrame('test', true);
      conn.socket.write(dataFrame.encode());
    });
  });

  it('Sends a CONTINUATION frame followed by a frame on a different stream', function (done) {
    var message = 'The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var headers = [
        [ ':method',    'GET'                      ],
        [ ':scheme',    'http'                     ],
        [ ':path',      '/'                        ],
        [ ':authority', helper.getHost()           ],
        [ 'x-dummy1',   helper.getDummyData(10000) ],
        [ 'x-dummy2',   helper.getDummyData(10000) ],
        [ 'x-dummy3',   helper.getDummyData(10000) ],
        [ 'x-dummy4',   helper.getDummyData(10000) ],
        [ 'x-dummy5',   helper.getDummyData(10000) ]
      ];

      var context = helper.createCompressionContext();
      var fragment = context.compress(headers);

      var headersFrame = helper.createHeadersFrame(fragment.slice(0, 16384), false, true);
      conn.socket.write(headersFrame.encode());

      var continuationFrame1 = helper.createContinuationFrame(fragment.slice(16384, 32767), false);
      conn.socket.write(continuationFrame1.encode());

      var continuationFrame2 = helper.createContinuationFrame(fragment.slice(32767), true);
      continuationFrame2.streamId = 3;
      conn.socket.write(continuationFrame2.encode());
    });
  });

  it('Sends a CONTINUATION frame with the stream identifier that is not 0x0', function (done) {
    var message = 'The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_PROTOCOL_ERROR, message);
        done();
      });

      var headers = [
        [ ':method',    'GET'                      ],
        [ ':scheme',    'http'                     ],
        [ ':path',      '/'                        ],
        [ ':authority', helper.getHost()           ],
        [ 'x-dummy1',   helper.getDummyData(10000) ],
        [ 'x-dummy2',   helper.getDummyData(10000) ]
      ];

      var context = helper.createCompressionContext();
      var fragment = context.compress(headers);

      var headersFrame = helper.createHeadersFrame(fragment.slice(0, 16384), false, true);
      conn.socket.write(headersFrame.encode());

      var continuationFrame = helper.createContinuationFrame(fragment.slice(16384), true);
      continuationFrame.streamId = 0x0;
      conn.socket.write(continuationFrame.encode());
    });
  });

  it('Sends a CONTINUATION frame after the frame other than HEADERS, PUSH_PROMISE or CONTINUATION', function (done) {
    var message = 'The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

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

      var headersFrame = helper.createHeadersFrame(fragment, true, false);
      conn.socket.write(headersFrame.encode());

      var dataFrame = helper.createDataFrame('test', true);
      conn.socket.write(dataFrame.encode());

      var continuationFrame = helper.createContinuationFrame(fragment, true);
      conn.socket.write(continuationFrame.encode());
    });
  });
});
