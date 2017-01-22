package http2

import (
	"fmt"

	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func SettingsSynchronization() *spec.TestGroup {
	tg := NewTestGroup("6.5.3", "Settings Synchronization")

	// The values in the SETTINGS frame MUST be processed in the order
	// they appear, with no other frame processing between values.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends multiple values of SETTINGS_INITIAL_WINDOW_SIZE",
		Requirement: "The endpoint MUST process the values in the settings in the order they apper.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			// Skip this test case when the length of data is 0.
			dataLen, err := spec.ServerDataLength(c)
			if err != nil {
				return err
			}
			if dataLen < 1 {
				return spec.ErrSkipped
			}

			err = conn.Handshake()
			if err != nil {
				return err
			}

			settings := []http2.Setting{
				http2.Setting{
					ID:  http2.SettingInitialWindowSize,
					Val: 100,
				},
				http2.Setting{
					ID:  http2.SettingInitialWindowSize,
					Val: 1,
				},
			}
			conn.WriteSettings(settings...)

			err = spec.VerifySettingsFrameWithAck(conn)
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			actual, passed := conn.WaitEventByType(spec.EventDataFrame)
			switch event := actual.(type) {
			case spec.DataFrameEvent:
				passed = (event.Header().Length == 1)
			default:
				passed = false
			}

			if !passed {
				expected := []string{
					fmt.Sprintf("DATA Frame (length:1, flags:0x00, stream_id:%d)", streamID),
				}

				return &spec.TestError{
					Expected: expected,
					Actual:   actual.String(),
				}
			}

			return nil
		},
	})

	// Once all values have been processed, the recipient MUST
	// immediately emit a SETTINGS frame with the ACK flag set.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a SETTINGS frame without ACK flag",
		Requirement: "The endpoint MUST immediately emit a SETTINGS frame with the ACK flag set.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

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
