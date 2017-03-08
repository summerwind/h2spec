package client

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Headers() *spec.ClientTestGroup {
	tg := NewTestGroup("6.2", "HEADERS")

	// END_HEADERS (0x4):
	// A HEADERS frame without the END_HEADERS flag set MUST be
	// followed by a CONTINUATION frame for the same stream.
	// A receiver MUST treat the receipt of any other type of frame
	// or a frame on a different stream as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	//
	// Note: This test case is duplicated with 4.3.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a HEADERS frame without the END_HEADERS flag, and a PRIORITY frame",
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

			pp := http2.PriorityParam{
				StreamDep: 0,
				Exclusive: false,
				Weight:    255,
			}
			conn.WritePriority(req.StreamID, pp)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// HEADERS frames MUST be associated with a stream. If a HEADERS
	// frame is received whose stream identifier field is 0x0, the
	// recipient MUST respond with a connection error (Section 5.4.1)
	// of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a HEADERS frame with 0x0 stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)

			hp := http2.HeadersFrameParam{
				StreamID:      0,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// The HEADERS frame can include padding. Padding fields and flags
	// are identical to those defined for DATA frames (Section 6.1).
	// Padding that exceeds the size remaining for the header block
	// fragment MUST be treated as a PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a HEADERS frame with invalid pad length",
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

			// HEADERS frame:
			// frame length: 16, pad length: 17
			var flags http2.Flags
			flags |= http2.FlagHeadersPadded
			payload := append([]byte("\x11"), conn.EncodeHeaders(headers)...)
			conn.WriteRawFrame(http2.FrameHeaders, flags, req.StreamID, payload)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
