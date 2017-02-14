package client

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func ExtendingHTTP2() *spec.ClientTestGroup {
	tg := NewTestGroup("5.5", "Extending HTTP/2")

	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends an unknown extension frame",
		Requirement: "The endpoint MUST ignore unknown or unsupported values in all extensible protocol elements.",
		Run: func(c *config.ClientSpecConfig, conn *spec.Conn, req *spec.Request) error {

			// UNKONWN Frame:
			// Length: 8, Type: 255, Flags: 0, R: 0, StreamID: 0
			conn.Send([]byte("\x00\x00\x08\x16\x00\x00\x00\x00\x00"))
			conn.Send([]byte("\x00\x00\x00\x00\x00\x00\x00\x00"))

			defer conn.WriteSuccessResponse(req.StreamID, c)

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	return tg
}
