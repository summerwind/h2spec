package generic

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Headers() *spec.TestGroup {
	tg := NewTestGroup("3.2", "HEADERS")

	// RFC7540, 6.2:
	// The HEADERS frame (type=0x1) is used to open a stream
	// (Section 5.1), and additionally carries a header block fragment.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame",
		Requirement: "The endpoint MUST accept HEADERS frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC7540, 6.2:
	// The HEADERS frame can include padding. Padding fields and flags
	// are identical to those defined for DATA frames (Section 6.1).
	// Padding that exceeds the size remaining for the header block
	// fragment MUST be treated as a PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with padding",
		Requirement: "The endpoint MUST accept HEADERS frame with padding.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
				PadLength:     8,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC7540, 6.2:
	// Prioritization information in a HEADERS frame is logically
	// equivalent to a separate PRIORITY frame, but inclusion in
	// HEADERS avoids the potential for churn in stream prioritization
	// when new streams are created. Prioritization fields in HEADERS
	// frames subsequent to the first on a stream reprioritize the
	// stream (Section 5.3.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with priority",
		Requirement: "The endpoint MUST accept HEADERS frame with priority.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			pp := http2.PriorityParam{
				StreamDep: 0,
				Exclusive: false,
				Weight:    255,
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
				Priority:      pp,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	return tg
}
