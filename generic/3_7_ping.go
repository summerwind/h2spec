package generic

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Ping() *spec.TestGroup {
	tg := NewTestGroup("3.7", "PING")

	// RFC7540, 6.7:
	// The PING frame (type=0x6) is a mechanism for measuring a minimal
	// round-trip time from the sender, as well as determining whether
	// an idle connection is still functional. PING frames can be sent
	// from any endpoint.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a PING frame",
		Requirement: "The endpoint MUST accept PING frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			data := [8]byte{'h', '2', 's', 'p', 'e', 'c'}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	return tg
}
