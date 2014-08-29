var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('4.2. Frame Size', function () {
  it('Sends too small size frame', function (done) {
    var message = 'The endpoint MUST send a FRAME_SIZE_ERROR error.';

    // WINDOW_UPDATE frame of length 8 octets
    var smallFrame = new Buffer([
      0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00,
      0x00
    ]);

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_FRAME_SIZE_ERROR, message);
        done();
      });

      conn.socket.write(smallFrame);
    });
  });

  it('Sends too large size frame', function (done) {
    var message = 'The endpoint MUST send a FRAME_SIZE_ERROR error.';

    // PING frame of length 16 octets
    var largeFrame = new Buffer([
      0x00, 0x00, 0x0f, 0x06, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00
    ]);

    var conn = helper.createConnection(function () {
      conn.once('error', function (err) {
        assert.equal(err.code, protocol.CODE_FRAME_SIZE_ERROR, message);
        done();
      });

      conn.socket.write(largeFrame);
    });
  });
});