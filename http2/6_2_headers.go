package http2

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Headers() *spec.TestGroup {
	tg := NewTestGroup("6.2", "HEADERS")

	// END_HEADERS (0x4):
	// A HEADERS frame without the END_HEADERS flag set MUST be
	// followed by a CONTINUATION frame for the same stream.
	// A receiver MUST treat the receipt of any other type of frame
	// or a frame on a different stream as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	//
	// Note: This test case is duplicated with 4.3.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame without the END_HEADERS flag, and a PRIORITY frame",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
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
			conn.WritePriority(streamID, pp)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// END_HEADERS (0x4):
	// A HEADERS frame without the END_HEADERS flag set MUST be
	// followed by a CONTINUATION frame for the same stream.
	// A receiver MUST treat the receipt of any other type of frame
	// or a frame on a different stream as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	//
	// Note: This test case is duplicated with 4.3.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame to another stream while sending a HEADERS frame",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)

			hp1 := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp1)

			hp2 := http2.HeadersFrameParam{
				StreamID:      streamID + 2,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp2)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// HEADERS frames MUST be associated with a stream. If a HEADERS
	// frame is received whose stream identifier field is 0x0, the
	// recipient MUST respond with a connection error (Section 5.4.1)
	// of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with 0x0 stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)

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
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with invalid pad length",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)

			fh := []byte("\x00\x00\x00\x01\x0d\x00\x00\x00\x01")
			fh[2] = byte(len(blockFragment) + 1)

			// HEADERS frame:
			// frame length: 16, pad length: 17
			conn.Send(fh)
			conn.Send([]byte{byte(len(blockFragment) + 2)})
			conn.Send(blockFragment)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
