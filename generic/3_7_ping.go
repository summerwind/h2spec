package generic

import (
	"fmt"
	"reflect"

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

			actual, passed := conn.WaitEventByType(spec.EventPingFrame)
			switch event := actual.(type) {
			case spec.PingFrameEvent:
				passed = (event.IsAck() && reflect.DeepEqual(event.Data, data))
			default:
				passed = false
			}

			if !passed {
				var actualStr string

				expected := []string{
					fmt.Sprintf("PING Frame (length:8, flags:0x01, stream_id:0, opaque_data:%s)", data),
				}

				f, ok := actual.(spec.PingFrameEvent)
				if ok {
					header := f.Header()
					actualStr = fmt.Sprintf(
						"PING Frame ((length:%d, flags:0x%02x, stream_id:%d, opaque_data: %s)",
						header.Type,
						header.Length,
						header.Flags,
						header.StreamID,
						f.Data,
					)
				} else {
					actualStr = actual.String()
				}

				return &spec.TestError{
					Expected: expected,
					Actual:   actualStr,
				}
			}

			return nil
		},
	})

	return tg
}
