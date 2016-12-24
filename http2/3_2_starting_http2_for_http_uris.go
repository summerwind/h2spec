package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

// The first HTTP/2 frame sent by the server MUST be a server connection
// preface (Section 3.5) consisting of a SETTINGS frame (Section 6.5).
func StartingHTTP2ForHTTPURIs() *spec.TestGroup {
	tg := NewTestGroup("3.2", "Starting HTTP/2 for \"http\" URIs")

	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends client connection preface",
		Requirement: "The first HTTP/2 frame sent by the server MUST be a server connection preface.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			setting := http2.Setting{
				ID:  http2.SettingInitialWindowSize,
				Val: spec.DefaultWindowSize,
			}

			conn.Send("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")
			conn.WriteSettings(setting)

			return spec.VerifyFrameType(conn, http2.FrameSettings)
		},
	})

	return tg
}
