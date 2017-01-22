package http2

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Continuation() *spec.TestGroup {
	tg := NewTestGroup("6.10", "CONTINUATION")

	// The CONTINUATION frame (type=0x9) is used to continue a sequence
	// of header block fragments (Section 4.3). Any number of
	// CONTINUATION frames can be sent, as long as the preceding frame
	// is on the same stream and is a HEADERS, PUSH_PROMISE,
	// or CONTINUATION frame without the END_HEADERS flag set.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends multiple CONTINUATION frames preceded by a HEADERS frame",
		Requirement: "The endpoint must accept the frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, false, conn.EncodeHeaders(dummyHeaders))
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// END_HEADERS (0x4):
	// If the END_HEADERS bit is not set, this frame MUST be followed
	// by another CONTINUATION frame. A receiver MUST treat the receipt
	// of any other type of frame or a frame on a different stream as
	// a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a CONTINUATION frame followed by any frame other than CONTINUATION",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, false, conn.EncodeHeaders(dummyHeaders))
			conn.WriteData(streamID, true, []byte("test"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// CONTINUATION frames MUST be associated with a stream. If a
	// CONTINUATION frame is received whose stream identifier field is
	// 0x0, the recipient MUST respond with a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a CONTINUATION frame with 0x0 stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(0, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A CONTINUATION frame MUST be preceded by a HEADERS, PUSH_PROMISE
	// or CONTINUATION frame without the END_HEADERS flag set.
	// A recipient that observes violation of this rule MUST respond
	// with a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a CONTINUATION frame preceded by a HEADERS frame with END_HEADERS flag",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A CONTINUATION frame MUST be preceded by a HEADERS, PUSH_PROMISE
	// or CONTINUATION frame without the END_HEADERS flag set.
	// A recipient that observes violation of this rule MUST respond
	// with a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a CONTINUATION frame preceded by a CONTINUATION frame with END_HEADERS flag",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A CONTINUATION frame MUST be preceded by a HEADERS, PUSH_PROMISE
	// or CONTINUATION frame without the END_HEADERS flag set.
	// A recipient that observes violation of this rule MUST respond
	// with a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a CONTINUATION frame preceded by a DATA frame",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)
			conn.WriteData(streamID, true, []byte("test"))

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, false, conn.EncodeHeaders(dummyHeaders))
			conn.WriteContinuation(0, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
