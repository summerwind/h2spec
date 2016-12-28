package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func ConnectionSpecificHeaderFields() *spec.TestGroup {
	tg := NewTestGroup("8.1.2.2", "Connection-Specific Header Fields")

	// An endpoint MUST NOT generate an HTTP/2 message containing
	// connection-specific header fields; any message containing
	// connection-specific header fields MUST be treated as
	// malformed (Section 8.1.2.6).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that contains the connection-specific header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField("connection", "keep-alive"))

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

	// The only exception to this is the TE header field, which MAY be
	// present in an HTTP/2 request; when it is, it MUST NOT contain
	// any value other than "trailers".
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that contains the TE header field with any value other than \"trailers\"",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField("trailers", "test"))
			headers = append(headers, spec.HeaderField("te", "trailers, deflate"))

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

	return tg
}
