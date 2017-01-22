package generic

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Continuation() *spec.TestGroup {
	tg := NewTestGroup("3.10", "CONTINUATION")

	// RFC7540, 6.10:
	// The CONTINUATION frame (type=0x9) is used to continue a sequence
	// of header block fragments (Section 4.3). Any number of
	// CONTINUATION frames can be sent, as long as the preceding frame
	// is on the same stream and is a HEADERS, PUSH_PROMISE, or
	// CONTINUATION frame without the END_HEADERS flag set.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a CONTINUATION frame",
		Requirement: "The endpoint MUST accept CONTINUATION frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headerBlock := conn.EncodeHeaders(headers)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: headerBlock[:5],
			}

			conn.WriteHeaders(hp)
			conn.WriteContinuation(streamID, true, headerBlock[5:])

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC7540, 6.10:
	// The CONTINUATION frame (type=0x9) is used to continue a sequence
	// of header block fragments (Section 4.3). Any number of
	// CONTINUATION frames can be sent, as long as the preceding frame
	// is on the same stream and is a HEADERS, PUSH_PROMISE, or
	// CONTINUATION frame without the END_HEADERS flag set.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends multiple CONTINUATION frames",
		Requirement: "The endpoint MUST accept multiple CONTINUATION frames.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headerBlock := conn.EncodeHeaders(headers)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    false,
				BlockFragment: headerBlock[:5],
			}

			conn.WriteHeaders(hp)
			conn.WriteContinuation(streamID, false, headerBlock[5:10])
			conn.WriteContinuation(streamID, true, headerBlock[10:])

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	return tg
}
