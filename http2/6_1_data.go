package http2

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Data() *spec.TestGroup {
	tg := NewTestGroup("6.1", "DATA")

	// DATA frames MUST be associated with a stream. If a DATA frame is
	// received whose stream identifier field is 0x0, the recipient
	// MUST respond with a connection error (Section 5.4.1) of type
	// PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a DATA frame with 0x0 stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteData(0, true, []byte("test"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// If a DATA frame is received whose stream is not in "open" or
	// "half-closed (local)" state, the recipient MUST respond with
	// a stream error (Section 5.4.2) of type STREAM_CLOSED.
	//
	// Note: This test case is duplicated with 5.1.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a DATA frame on the stream that is not in \"open\" or \"half-closed (local)\" state",
		Requirement: "The endpoint MUST respond with a stream error of type STREAM_CLOSED.",
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
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)
			conn.WriteData(streamID, true, []byte("test"))

			return spec.VerifyStreamError(conn, http2.ErrCodeStreamClosed)
		},
	})

	// If the length of the padding is the length of the frame payload
	// or greater, the recipient MUST treat this as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a DATA frame with invalid pad length",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"
			headers = append(headers, spec.HeaderField("content-length", "4"))

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)

			// DATA frame:
			// frame length: 5, pad length: 6
			conn.Send([]byte("\x00\x00\x05\x00\x09\x00\x00\x00\x01"))
			conn.Send([]byte("\x06\x54\x65\x73\x74"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
