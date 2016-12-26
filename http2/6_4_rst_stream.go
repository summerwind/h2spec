package http2

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func RSTStream() *spec.TestGroup {
	tg := NewTestGroup("6.4", "RST_STREAM")

	// RST_STREAM frames MUST be associated with a stream.  If a
	// RST_STREAM frame is received with a stream identifier of 0x0,
	// the recipient MUST treat this as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a RST_STREAM frame with 0x0 stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteRSTStream(0, http2.ErrCodeCancel)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// RST_STREAM frames MUST NOT be sent for a stream in the "idle"
	// state. If a RST_STREAM frame identifying an idle stream is
	// received, the recipient MUST treat this as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a RST_STREAM frame on a idle stream",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteRSTStream(1, http2.ErrCodeCancel)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A RST_STREAM frame with a length other than 4 octets MUST be
	// treated as a connection error (Section 5.4.1) of type
	// FRAME_SIZE_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a RST_STREAM frame with a length other than 4 octets",
		Requirement: "The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
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

			// RST_STREAM frame:
			// length: 3, flags: 0x0, stream_id: 0x01
			conn.Send([]byte("\x00\x00\x03\x03\x00\x00\x00\x00\x01"))
			conn.Send([]byte("\x00\x00\x00"))

			return spec.VerifyStreamError(conn, http2.ErrCodeFrameSize)
		},
	})

	return tg
}
