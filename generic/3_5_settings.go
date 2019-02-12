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
				http2.Setting{
					ID:  http2.SettingHeaderTableSize,
					Val: 4096,
				},
				http2.Setting{
					ID:  http2.SettingEnablePush,
					Val: 1,
				},
				http2.Setting{
					ID:  http2.SettingMaxConcurrentStreams,
					Val: 100,
				},
				http2.Setting{
					ID:  http2.SettingInitialWindowSize,
					Val: 65535,
				},
				http2.Setting{
					ID:  http2.SettingMaxFrameSize,
					Val: 16384,
				},
				http2.Setting{
					ID:  http2.SettingMaxHeaderListSize,
					Val: 100,
				},
			}
			conn.WriteSettings(settings...)

			return spec.VerifySettingsFrameWithAck(conn)
		},
	})

	return tg
}
