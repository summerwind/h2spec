package client

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func SettingsSynchronization() *spec.ClientTestGroup {
	tg := NewTestGroup("6.5.3", "Settings Synchronization")

	// Once all values have been processed, the recipient MUST
	// immediately emit a SETTINGS frame with the ACK flag set.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a SETTINGS frame without ACK flag",
		Requirement: "The endpoint MUST immediately emit a SETTINGS frame with the ACK flag set.",
		Run: func(c *config.ClientSpecConfig, conn *spec.Conn, req *spec.Request) error {
			setting := http2.Setting{
				ID:  http2.SettingEnablePush,
				Val: 0,
			}
			conn.WriteSettings(setting)

			return spec.VerifySettingsFrameWithAck(conn)
		},
	})

	return tg
}
