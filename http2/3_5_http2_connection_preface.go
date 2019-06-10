package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func HTTP2ConnectionPreface() *spec.TestGroup {
	tg := NewTestGroup("3.5", "HTTP/2 Connection Preface")

	// The server connection preface consists of a potentially empty
	// SETTINGS frame (Section 6.5) that MUST be the first frame
	// the server sends in the HTTP/2 connection.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends client connection preface",
		Requirement: "The server connection preface MUST be the first frame the server sends in the HTTP/2 connection.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var passed bool

			setting := http2.Setting{
				ID:  http2.SettingInitialWindowSize,
				Val: spec.DefaultWindowSize,
			}

			conn.Send([]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"))
			conn.WriteSettings(setting)

			actual := conn.WaitEvent()
			switch event := actual.(type) {
			case spec.SettingsFrameEvent:
				passed = !event.IsAck()
			default:
				passed = false
			}

			if !passed {
				expected := []string{
					"SETTINGS Frame (flags:0x00)",
				}

				return &spec.TestError{
					Expected: expected,
					Actual:   actual.String(),
				}
			}

			return nil
		},
	})

	// Clients and servers MUST treat an invalid connection preface as
	// a connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends invalid connection preface",
		Requirement: "The endpoint MUST terminate the TCP connection.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Send([]byte("INVALID CONNECTION PREFACE\r\n\r\n"))
			if err != nil {
				return err
			}

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
