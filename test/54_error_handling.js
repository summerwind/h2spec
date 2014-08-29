var assert = require('assert');

var helper = require('./helper'),
    protocol = helper.protocol;

describe('5.4. Error Handling', function () {
  describe('5.4.1. Connection Error Handling', function () {
    it('Receives a GOAWAY frame', function (done) {
      var message = 'After sending the GOAWAY frame, the endpoint MUST close the TCP connection.';

      var conn = helper.createConnection(function () {
        conn.once('error', function () {});
        conn.once('close', function (hasError) {
          assert(hasError, message);
          done();
        });

        var frame = helper.createDataFrame('test', true);
        conn.socket.write(frame.encode());
      });
    });
  });
})
