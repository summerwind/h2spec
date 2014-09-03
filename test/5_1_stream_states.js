var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('5.1. Stream States', function () {
  describe('5.1.1. Stream Identifiers', function () {
    var message = 'The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.';

    it('Sends even-numbered stream identifier', function (done) {
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

        var frame = helper.createHeadersFrame(fragment, true, true);
        frame.streamId = 2;
        conn.socket.write(frame.encode());
      });
    });

    it('Sends stream identifier that is numerically smaller than previous', function (done) {
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

        var frame = helper.createHeadersFrame(fragment, true, true);
        frame.streamId = 3;
        conn.socket.write(frame.encode());

        var frame = helper.createHeadersFrame(fragment, true, true);
        frame.streamId = 1;
        conn.socket.write(frame.encode());
      });
    });
  });
});