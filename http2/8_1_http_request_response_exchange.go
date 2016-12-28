package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func HTTPRequestResponseExchange() *spec.TestGroup {
	tg := NewTestGroup("8.1", "HTTP Request/Response Exchange")

	// An endpoint that receives a HEADERS frame without the
	// END_STREAM flag set after receiving a final (non-informational)
	// status code MUST treat the corresponding request or response
	// as malformed (Section 8.1.2.6).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a second HEADERS frame without the END_STREAM flag",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"

			hp1 := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp1)
			conn.WriteData(streamID, false, []byte("test"))

			trailers := []hpack.HeaderField{
				spec.HeaderField("x-test", "ok"),
			}

			hp2 := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(trailers),
			}

			conn.WriteHeaders(hp2)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	tg.AddTestGroup(HTTPHeaderFields())

	return tg
}
