package generic

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Settings() *spec.TestGroup {
	tg := NewTestGroup("3.5", "SETTINGS")

	// RFC7540, 6.5:
	// he SETTINGS frame (type=0x4) conveys configuration parameters
	// that affect how endpoints communicate, such as preferences and
	// constraints on peer behavior. The SETTINGS frame is also used
	// to acknowledge the receipt of those parameters. Individually,
	// a SETTINGS parameter can also be referred to as a "setting".
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a SETTINGS frame",
		Requirement: "The endpoint MUST accept SETTINGS frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			settings := []http2.Setting{
				http2.Setting{http2.SettingHeaderTableSize, 4096},
				http2.Setting{http2.SettingEnablePush, 1},
				http2.Setting{http2.SettingMaxConcurrentStreams, 100},
				http2.Setting{http2.SettingInitialWindowSize, 65535},
				http2.Setting{http2.SettingMaxFrameSize, 16384},
				http2.Setting{http2.SettingMaxHeaderListSize, 100},
			}
			conn.WriteSettings(settings...)

			actual, passed := conn.WaitEventByType(spec.EventSettingsFrame)
			switch event := actual.(type) {
			case spec.SettingsFrameEvent:
				passed = event.IsAck()
			default:
				passed = false
			}

			if !passed {
				expected := []string{
					"SETTINGS Frame (flags:0x01)",
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
