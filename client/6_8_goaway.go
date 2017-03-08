package client

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func GoAway() *spec.ClientTestGroup {
	tg := NewTestGroup("6.8", "GOAWAY")

	// An endpoint MUST treat a GOAWAY frame with a stream identifier
	// other than 0x0 as a connection error (Section 5.4.1) of type
	// PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a GOAWAY frame with a stream identifier other than 0x0",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// GOAWAY frame:
			// length: 8, flags: 0x0, stream_id: 1
			conn.Send([]byte("\x00\x00\x08\x07\x00\x00\x00\x00\x01"))
			conn.Send([]byte("\x00\x00\x00\x00\x00\x00\x00\x00"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
