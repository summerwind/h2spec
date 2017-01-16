package hpack

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func IndexedHeaderFieldRepresentation() *spec.TestGroup {
	tg := NewTestGroup("6.1", "Indexed Header Field Representation")

	// The index value of 0 is not used.  It MUST be treated as a decoding
	// error if found in an indexed header field representation.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a indexed header field representation with index 0",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Indexed Header Field Representation
			rep := []byte("\x80")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(blockFragment, rep...)

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
