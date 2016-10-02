package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StreamConcurrency() *spec.TestGroup {
	tg := NewTestGroup("5.1.2", "Stream Concurrency")

	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends HEADERS frames that causes their advertised concurrent stream limit to be exceeded",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR or REFUSED_STREAM.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Skip this when SETTINGS_MAX_CONCURRENT_STREAMS is unlimited.
			maxStreams, ok := conn.Settings[http2.SettingMaxConcurrentStreams]
			if !ok {
				return spec.ErrSkipped
			}

			// Set INITIAL_WINDOW_SIZE to zero to prevent the peer from closing the stream
			settings := http2.Setting{http2.SettingInitialWindowSize, 0}
			conn.WriteSettings(settings)

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)

			var streamID uint32 = 1
			for i := 0; i <= int(maxStreams); i++ {
				hp := http2.HeadersFrameParam{
					StreamID:      streamID,
					EndStream:     true,
					EndHeaders:    true,
					BlockFragment: blockFragment,
				}
				conn.WriteHeaders(hp)
				streamID += 2
			}

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol, http2.ErrCodeRefusedStream)
		},
	})

	return tg
}
