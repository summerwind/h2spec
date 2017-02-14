package client

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func StartingHTTP2() *spec.ClientTestGroup {
	tg := NewTestGroup("1", "Starting HTTP/2")

	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a client connection preface",
		Requirement: "The endpoint MUST accept client connection preface.",
		Run: func(c *config.ClientSpecConfig, conn *spec.Conn, req *spec.Request) error {
			conn.WriteSuccessResponse(req.StreamID, c)

			return nil
		},
	})

	return tg
}
