package http2

import (
	"fmt"

	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func ConnectionErrorHandling() *spec.TestGroup {
	tg := NewTestGroup("5.4.1", "Connection Error Handling")

	// After sending the GOAWAY frame for an error condition,
	// the endpoint MUST close the TCP connection.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends an invalid PING frame for connection close",
		Requirement: "The endpoint MUST close the TCP connection",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// PING frame with invalid stream ID
			conn.Send([]byte("\x00\x00\x08\x06\x00\x00\x00\x00\x03"))
			conn.Send([]byte("\x00\x00\x00\x00\x00\x00\x00\x00"))

			return spec.VerifyConnectionClose(conn)
		},
	})

	// An endpoint that encounters a connection error SHOULD first send
	// a GOAWAY frame (Section 6.8) with the stream identifier of the last
	// stream that it successfully received from its peer.
	tg.AddTestCase(&spec.TestCase{
		Strict:      true,
		Desc:        "Sends an invalid PING frame to receive GOAWAY frame",
		Requirement: "An endpoint that encounters a connection error SHOULD first send a GOAWAY frame",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// PING frame with invalid stream ID
			conn.Send([]byte("\x00\x00\x08\x06\x00\x00\x00\x00\x03"))
			conn.Send([]byte("\x00\x00\x00\x00\x00\x00\x00\x00"))

			actual, passed := conn.WaitEventByType(spec.EventGoAwayFrame)
			if !passed {
				return &spec.TestError{
					Expected: []string{
						fmt.Sprintf(spec.ExpectedGoAwayFrame, http2.ErrCodeProtocol),
					},
					Actual: actual.String(),
				}
			}

			return nil
		},
	})

	return tg
}
