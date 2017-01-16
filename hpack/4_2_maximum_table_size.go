package hpack

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func MaximumTableSize() *spec.TestGroup {
	tg := NewTestGroup("4.2", "Maximum Table Size")

	// A change in the maximum size of the dynamic table is signaled
	// via a dynamic table size update (see Section 6.3). This dynamic
	// table size update MUST occur at the beginning of the first
	// header block following the change to the dynamic table size.
	// In HTTP/2, this follows a settings acknowledgment (see Section
	// 6.5.3 of [HTTP2]).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a dynamic table size update at the end of header block",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Dynamic table size update with value 1
			tableSizeUpdate := []byte("\x21")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(blockFragment, tableSizeUpdate...)

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
