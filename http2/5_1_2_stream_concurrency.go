package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StreamConcurrency() *spec.TestGroup {
	tg := NewTestGroup("5.1.2", "Stream Concurrency")

	// An endpoint that receives a HEADERS frame that causes
	// its advertised concurrent stream limit to be exceeded
	// MUST treat this as a stream error (Section 5.4.2) of
	// type PROTOCOL_ERROR or REFUSED_STREAM.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends HEADERS frames that causes their advertised concurrent stream limit to be exceeded",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR or REFUSED_STREAM.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Skip this test case when SETTINGS_MAX_CONCURRENT_STREAMS
			// is unlimited.
			maxStreams, ok := conn.Settings[http2.SettingMaxConcurrentStreams]
			if !ok {
				return spec.ErrSkipped
			}

			// Set INITIAL_WINDOW_SIZE to zero to prevent the peer from
			// closing the stream.
			settings := http2.Setting{
				ID:  http2.SettingInitialWindowSize,
				Val: 0,
			}
			conn.WriteSettings(settings)

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)

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

			codes := []http2.ErrCode{
				http2.ErrCodeProtocol,
				http2.ErrCodeRefusedStream,
			}
			return spec.VerifyStreamError(conn, codes...)
		},
	})

	return tg
}
