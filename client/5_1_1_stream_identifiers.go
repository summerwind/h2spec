package client

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StreamIdentifiers() *spec.ClientTestGroup {
	tg := NewTestGroup("5.1.1", "Stream Identifiers")

	// An endpoint that receives an unexpected stream identifier
	// MUST respond with a connection error (Section 5.4.1) of
	// type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends an unexpected stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      101,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
