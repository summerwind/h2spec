package http2

import (
	"fmt"

	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func ConnectionErrorHandling() *spec.TestGroup {
	tg := NewTestGroup("5.4.1", "Connection Error Handling")

	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends invalid frame for connection close",
		Requirement: "The endpoint MUST close the TCP connection",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// PING frame with invalid stream ID
			conn.Send("\x00\x00\x08\x06\x00\x00\x00\x00\x03")
			conn.Send("\x00\x00\x00\x00\x00\x00\x00\x00")

			return spec.VerifyConnectionClose(conn)
		},
	})

	tg.AddTestCase(&spec.TestCase{
		Strict:      true,
		Desc:        "Sends invalid frame for GOAWAY frame",
		Requirement: "An endpoint that encounters a connection error SHOULD first send a GOAWAY frame",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var actual spec.Event

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// PING frame with invalid stream ID
			conn.Send("\x00\x00\x08\x06\x00\x00\x00\x00\x03")
			conn.Send("\x00\x00\x00\x00\x00\x00\x00\x00")

			passed := false
			for !conn.Closed {
				actual = conn.WaitEvent()
				_, passed = actual.(spec.EventGoAwayFrame)
				if passed {
					break
				}
			}

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
