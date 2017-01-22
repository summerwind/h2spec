package http2

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func DefinedSETTINGSParameters() *spec.TestGroup {
	tg := NewTestGroup("6.5.2", "Defined SETTINGS Parameters")

	// SETTINGS_ENABLE_PUSH (0x2):
	// The initial value is 1, which indicates that server push is
	// permitted. Any value other than 0 or 1 MUST be treated as a
	// connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "SETTINGS_ENABLE_PUSH (0x2): Sends the value other than 0 or 1",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			setting := http2.Setting{
				ID:  http2.SettingEnablePush,
				Val: 2,
			}
			conn.WriteSettings(setting)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// SETTINGS_INITIAL_WINDOW_SIZE (0x4):
	// Values above the maximum flow-control window size of 2^31-1
	// MUST be treated as a connection error (Section 5.4.1) of
	// type FLOW_CONTROL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "SETTINGS_INITIAL_WINDOW_SIZE (0x4): Sends the value above the maximum flow control window size",
		Requirement: "The endpoint MUST treat this as a connection error of type FLOW_CONTROL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			setting := http2.Setting{
				ID:  http2.SettingInitialWindowSize,
				Val: 2147483648,
			}
			conn.WriteSettings(setting)

			return spec.VerifyConnectionError(conn, http2.ErrCodeFlowControl)
		},
	})

	// SETTINGS_MAX_FRAME_SIZE (0x5):
	// The initial value is 2^14 (16,384) octets. The value advertised
	// by an endpoint MUST be between this initial value and the
	// maximum allowed frame size (2^24-1 or 16,777,215 octets),
	// inclusive. Values outside this range MUST be treated as a
	// connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "SETTINGS_MAX_FRAME_SIZE (0x5): Sends the value below the initial value",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			setting := http2.Setting{
				ID:  http2.SettingMaxFrameSize,
				Val: 16383,
			}
			conn.WriteSettings(setting)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// SETTINGS_MAX_FRAME_SIZE (0x5):
	// The initial value is 2^14 (16,384) octets. The value advertised
	// by an endpoint MUST be between this initial value and the
	// maximum allowed frame size (2^24-1 or 16,777,215 octets),
	// inclusive. Values outside this range MUST be treated as a
	// connection error (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "SETTINGS_MAX_FRAME_SIZE (0x5): Sends the value above the maximum allowed frame size",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			setting := http2.Setting{
				ID:  http2.SettingMaxFrameSize,
				Val: 16777216,
			}
			conn.WriteSettings(setting)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// An endpoint that receives a SETTINGS frame with any unknown
	// or unsupported identifier MUST ignore that setting.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a SETTINGS frame with unknown identifier",
		Requirement: "The endpoint MUST ignore that setting.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			setting := http2.Setting{
				ID:  0xFF,
				Val: 1,
			}
			conn.WriteSettings(setting)

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	return tg
}
