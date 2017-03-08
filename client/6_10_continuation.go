package client

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Continuation() *spec.ClientTestGroup {
	tg := NewTestGroup("6.10", "CONTINUATION")

	// The CONTINUATION frame (type=0x9) is used to continue a sequence
	// of header block fragments (Section 4.3). Any number of
	// CONTINUATION frames can be sent, as long as the preceding frame
	// is on the same stream and is a HEADERS, PUSH_PROMISE,
	// or CONTINUATION frame without the END_HEADERS flag set.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends multiple CONTINUATION frames preceded by a HEADERS frame",
		Requirement: "The endpoint must accept the frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, false, conn.EncodeHeaders(dummyHeaders))
			conn.WriteContinuation(req.StreamID, true, conn.EncodeHeaders(dummyHeaders))

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	// END_HEADERS (0x4):
	// If the END_HEADERS bit is not set, this frame MUST be followed
	// by another CONTINUATION frame. A receiver MUST treat the receipt
	// of any other type of frame or a frame on a different stream as
	// a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a CONTINUATION frame followed by any frame other than CONTINUATION",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)

			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     false,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, false, conn.EncodeHeaders(dummyHeaders))
			conn.WriteData(req.StreamID, true, []byte("test"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// CONTINUATION frames MUST be associated with a stream. If a
	// CONTINUATION frame is received whose stream identifier field is
	// 0x0, the recipient MUST respond with a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a CONTINUATION frame with 0x0 stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(0, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A CONTINUATION frame MUST be preceded by a HEADERS, PUSH_PROMISE
	// or CONTINUATION frame without the END_HEADERS flag set.
	// A recipient that observes violation of this rule MUST respond
	// with a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a CONTINUATION frame preceded by a HEADERS frame with END_HEADERS flag",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A CONTINUATION frame MUST be preceded by a HEADERS, PUSH_PROMISE
	// or CONTINUATION frame without the END_HEADERS flag set.
	// A recipient that observes violation of this rule MUST respond
	// with a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a CONTINUATION frame preceded by a CONTINUATION frame with END_HEADERS flag",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, true, conn.EncodeHeaders(dummyHeaders))
			conn.WriteContinuation(req.StreamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A CONTINUATION frame MUST be preceded by a HEADERS, PUSH_PROMISE
	// or CONTINUATION frame without the END_HEADERS flag set.
	// A recipient that observes violation of this rule MUST respond
	// with a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a CONTINUATION frame preceded by a DATA frame",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)

			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     false,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)
			conn.WriteData(req.StreamID, true, []byte("test"))

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, false, conn.EncodeHeaders(dummyHeaders))
			conn.WriteContinuation(0, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
