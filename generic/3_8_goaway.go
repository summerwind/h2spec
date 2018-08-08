package generic

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func GoAway() *spec.TestGroup {
	tg := NewTestGroup("3.8", "GOAWAY")

	// RFC7540, 6.8:
	// The GOAWAY frame (type=0x7) is used to initiate shutdown of a
	// connection or to signal serious error conditions. GOAWAY allows
	// an endpoint to gracefully stop accepting new streams while
	// still finishing processing of previously established streams.
	// This enables administrative actions, like server maintenance.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a GOAWAY frame",
		Requirement: "The endpoint MUST accept GOAWAY frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteGoAway(0, http2.ErrCodeNo, []byte("h2spec"))

			data := [8]byte{'h', '2', 's', 'p', 'e', 'c'}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameOrConnectionClose(conn, data)
		},
	})

	return tg
}
