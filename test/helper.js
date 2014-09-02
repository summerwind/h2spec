var net = require('net');

var Connection = require('../node_modules/sasazka/lib/connection'),
    framer = require('../node_modules/sasazka/lib/framer'),
    hpack = require('../node_modules/sasazka/lib/hpack'),
    protocol = require('../node_modules/sasazka/lib/protocol');

function getPort() {
  return process.env.H2CHECK_PORT;
}

function getHost() {
  return process.env.H2CHECK_HOST;
}

function createConnection(cb) {
  var port = getPort();
  var host = getHost();

  var socket = net.connect(port, host);
  var conn = new Connection(socket, {});
  conn.on('connect', function () {
    var frameHandler = function (frame) {
      if (frame.type === protocol.FRAME_TYPE_SETTINGS && frame.ack) {
        conn.removeListener('frame', frameHandler);
        setImmediate(cb);
      }
    };

    conn.on('frame', frameHandler);
  });

  return conn;
};

function createDataFrame(data, end) {
  var frame = framer.createDataFrame();
  frame.streamId = 1;
  frame.setData(data);

  if (end) {
    frame.endStream = true;
  }

  return frame;
};

function createHeadersFrame(fragment, endHeaders, endStream) {
  var frame = framer.createHeadersFrame();
  frame.streamId = 1;
  frame.setHeaderBlockFragment(fragment);

  if (endHeaders) {
    frame.endHeaders = true;
  }
  if (endStream) {
    frame.endStream = true;
  }

  return frame;
};

function createPriorityFrame(dependency, weight, exclusive) {
  var frame = framer.createPriorityFrame();
  frame.streamId = 1;
  frame.setPriority(dependency, weight, exclusive);

  return frame;
};

function createRstStreamFrame(errorCode) {
  var frame = framer.createRstStreamFrame();
  frame.streamId = 1;
  frame.setErrorCode(errorCode);

  return frame;
};

function createSettingsFrame() {
  var frame = framer.createSettingsFrame();
  frame.streamId = 0;

  return frame;
};

function createPingFrame(ack) {
  var frame = framer.createPingFrame();
  frame.streamId = 0;
  frame.ack = (ack === true);

  return frame;
};

function createGoawayFrame(lastStreamId, errorCode, debugData) {
  var frame = framer.createGoawayFrame();
  frame.streamId = 0;
  frame.setLastStreamId(lastStreamId);
  frame.setErrorCode(errorCode);

  if (debugData) {
    frame.setDebugData(debugData);
  }

  return frame;
};

function createWindowUpdateFrame(increment) {
  var frame = framer.createWindowUpdateFrame();
  frame.streamId = 0;
  frame.setWindowSizeIncrement(increment);

  return frame;
};

function createContinuationFrame(fragment, end) {
  var frame = framer.createContinuationFrame();
  frame.streamId = 1;
  frame.setHeaderBlockFragment(fragment);

  if (end) {
    frame.endHeaders = true;
  }

  return frame;
};

function createCompressionContext() {
  return hpack.createContext();
};

module.exports = {
  protocol: protocol,
  getPort: getPort,
  getHost: getHost,
  createConnection: createConnection,
  createDataFrame: createDataFrame,
  createHeadersFrame: createHeadersFrame,
  createPriorityFrame: createPriorityFrame,
  createRstStreamFrame: createRstStreamFrame,
  createSettingsFrame: createSettingsFrame,
  createPingFrame: createPingFrame,
  createGoawayFrame: createGoawayFrame,
  createWindowUpdateFrame: createWindowUpdateFrame,
  createContinuationFrame: createContinuationFrame,
  createCompressionContext: createCompressionContext
};
