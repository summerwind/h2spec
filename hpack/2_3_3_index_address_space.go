package hpack

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func IndexAddressSpace() *spec.TestGroup {
	tg := NewTestGroup("2.3.3", "Index Address Space")

	// Indices strictly greater than the sum of the lengths of both
	// tables MUST be treated as a decoding error.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a indexed header field representation with invalid index",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Indexed header field representation with index 70
			indexedRep := []byte("\xC6")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(blockFragment, indexedRep...)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: blockFragment,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeCompression)
		},
	})

	// Indices strictly greater than the sum of the lengths of both
	// tables MUST be treated as a decoding error.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field representation with invalid index",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal Header Field with Incremental Indexing (index=70 & value=empty)
			indexedRep := []byte("\x7F\x07\x00")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(blockFragment, indexedRep...)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: blockFragment,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeCompression)
		},
	})

	return tg
}
