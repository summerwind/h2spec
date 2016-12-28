package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func HTTPHeaderFields() *spec.TestGroup {
	tg := NewTestGroup("8.1.2", "HTTP Header Fields")

	// Just as in HTTP/1.x, header field names are strings of ASCII
	// characters that are compared in a case-insensitive fashion.
	// However, header field names MUST be converted to lowercase
	// prior to their encoding in HTTP/2. A request or response
	// containing uppercase header field names MUST be treated as
	// malformed (Section 8.1.2.6).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that contains the header field name in uppercase letters",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField("X-TEST", "ok"))

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	tg.AddTestGroup(PseudoHeaderFields())
	tg.AddTestGroup(ConnectionSpecificHeaderFields())
	tg.AddTestGroup(RequestPseudoHeaderFields())
	tg.AddTestGroup(MalformedRequestsAndResponses())

	return tg
}
