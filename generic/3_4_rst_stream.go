package generic

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func RSTStream() *spec.TestGroup {
	tg := NewTestGroup("3.4", "RST_STREAM")

	// RFC7540, 6.4:
	// The RST_STREAM frame (type=0x3) allows for immediate termination
	// of a stream. RST_STREAM is sent to request cancellation of a
	// stream or to indicate that an error condition has occurred.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a RST_STREAM frame",
		Requirement: "The endpoint MUST accept RST_STREAM frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			conn.WriteRSTStream(streamID, http2.ErrCodeCancel)

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	return tg
}
