package client

import (
	"errors"
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
		Run: func(c *config.Config, conn *spec.Conn) error {
			_, err := conn.ReadClientPreface()
			if err != nil {
				return err
			}

			_, ok := conn.WaitEventByType(spec.EventSettingsFrame)
			if !ok {
				return errors.New("First frame from client must be SETTINGS")
			}

			setting := http2.Setting{
				ID:  http2.SettingMaxConcurrentStreams,
				Val: 100,
			}
			conn.WriteSettings(setting)

			return spec.VerifySettingsFrameWithAck(conn)
		},
	})

	return tg
}
