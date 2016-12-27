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
		Desc:        "Sends multiple values of MAX_CONCURRENT_STREAMS",
		Requirement: "The endpoint MUST process the values in the settings in the order they apper.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1
			var actual spec.Event

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp1 := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp1)

			// Get the length of response body.
			resLen := -1
			for resLen == -1 {
				ev := conn.WaitEvent()

				switch event := ev.(type) {
				case spec.EventDataFrame:
					resLen = int(event.Header().Length)
				}
			}

			// Skip this test case when the length of response body is 0.
			if resLen < 1 {
				return spec.ErrSkipped
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

			err = spec.VerifyFrameType(conn, http2.FrameSettings)
			if err != nil {
				return err
			}

			streamID += 2
			hp2 := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp2)

			passed := false
			for !conn.Closed {
				ev := conn.WaitEvent()

				switch event := ev.(type) {
				case spec.EventDataFrame:
					actual = event
					passed = (event.Header().Length == 1)
				case spec.EventTimeout:
					if actual == nil {
						actual = event
					}
				default:
					actual = ev
				}

				if passed {
					break
				}
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
			var actual spec.Event

			err := conn.Handshake()
			if err != nil {
				return err
			}

			setting := http2.Setting{
				ID:  http2.SettingEnablePush,
				Val: 0,
			}
			conn.WriteSettings(setting)

			passed := false
			for !conn.Closed {
				ev := conn.WaitEvent()

				switch event := ev.(type) {
				case spec.EventSettingsFrame:
					actual = event
					passed = event.IsAck()
				case spec.EventTimeout:
					if actual == nil {
						actual = event
					}
				default:
					actual = ev
				}

				if passed {
					break
				}
			}

			if !passed {
				expected := []string{
					"SETTINGS Frame (length:0, flags:0x01, stream_id:0)",
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
