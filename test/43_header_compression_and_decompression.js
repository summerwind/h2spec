var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('4.3. Header Compression and Decompression', function () {
  it('Sends invalid header block fragment', function (done) {
    var message = 'The endpoint MUST terminate the connection with a connection error of type COMPRESSION_ERROR.';

    // HEADERS frame with invalid header blocks
    var frame = new Buffer([
      0x00, 0x00, 0x14, 0x01, 0x05, 0x00, 0x00, 0x00,
      0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00
    ]);

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_COMPRESSION_ERROR, message);
        done();
      });

      conn.socket.write(frame);
    });
  });
});