package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func HTTP2ConnectionPreface() *spec.TestGroup {
	tg := NewTestGroup("3.5", "HTTP/2 Connection Preface")

	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends invalid connection preface",
		Requirement: "The endpoint MUST terminate the TCP connection.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Send("INVALID CONNECTION PREFACE\r\n\r\n")
			if err != nil {
				return err
			}

			return spec.VerifyConnectionClose(conn)
		},
	})

	return tg
}
