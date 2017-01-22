package http2

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func ExtendingHTTP2() *spec.TestGroup {
	tg := NewTestGroup("5.5", "Extending HTTP/2")

	// Implementations MUST ignore unknown or unsupported values
	// in all extensible protocol elements.
	//
	// Note: This test case is duplicated with 4.1.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends an unknown extension frame",
		Requirement: "The endpoint MUST ignore unknown or unsupported values in all extensible protocol elements.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// UNKONWN Frame:
			// Length: 8, Type: 255, Flags: 0, R: 0, StreamID: 0
			conn.Send([]byte("\x00\x00\x08\x16\x00\x00\x00\x00\x00"))
			conn.Send([]byte("\x00\x00\x00\x00\x00\x00\x00\x00"))

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	// Extension frames that appear in the middle of a header block
	// (Section 4.3) are not permitted; these MUST be treated as
	// a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends an unknown extension frame in the middle of a header block",
		Requirement: "The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.",
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

			// UNKONWN Frame:
			// Length: 8, Type: 255, Flags: 0, R: 0, StreamID: 0
			conn.Send([]byte("\x00\x00\x08\x16\x00\x00\x00\x00\x00"))
			conn.Send([]byte("\x00\x00\x00\x00\x00\x00\x00\x00"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
