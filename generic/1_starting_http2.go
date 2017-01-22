package generic

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StartingHTTP2() *spec.TestGroup {
	tg := NewTestGroup("1", "Starting HTTP/2")

	// RFC7540, 3.2:
	// The first HTTP/2 frame sent by the server MUST be a server connection
	// preface (Section 3.5) consisting of a SETTINGS frame (Section 6.5).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a client connection preface",
		Requirement: "The endpoint MUST accept client connection preface.",
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

	return tg
}
