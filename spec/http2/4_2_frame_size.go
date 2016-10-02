package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func FrameSize() *spec.TestGroup {
	tg := NewTestGroup("4.2", "Frame Size")

	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends large size frame that exceeds the SETTINGS_MAX_FRAME_SIZE",
		Requirement: "The endpoint MUST send a FRAME_SIZE_ERROR error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			maxFrameSize, ok := conn.Settings[http2.SettingMaxFrameSize]
			if !ok {
				maxFrameSize = 18384
			}

			payload := spec.DummyBytes(int(maxFrameSize) + 1)

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)
			conn.WriteData(streamID, true, payload)

			return spec.VerifyStreamError(conn, http2.ErrCodeFrameSize)
		},
	})

	return tg
}
