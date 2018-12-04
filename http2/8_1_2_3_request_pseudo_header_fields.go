package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func RequestPseudoHeaderFields() *spec.TestGroup {
	tg := NewTestGroup("8.1.2.3", "Request Pseudo-Header Fields")

	// The ":path" pseudo-header field includes the path and query
	// parts of the target URI (the "path-absolute" production and
	// optionally a '?' character followed by the "query" production
	// (see Sections 3.3 and 3.4 of [RFC3986]). A request in asterisk
	// form includes the value '*' for the ":path" pseudo-header field.
	//
	// This pseudo-header field MUST NOT be empty for "http" or "https"
	// URIs; "http" or "https" URIs that do not contain a path
	// component MUST include a value of '/'.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with empty \":path\" pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[2].Value = ""

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

	// All HTTP/2 requests MUST include exactly one valid value for
	// the ":method", ":scheme", and ":path" pseudo-header fields,
	// unless it is a CONNECT request (Section 8.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that omits \":method\" pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = []hpack.HeaderField{
				headers[1], // :scheme
				headers[2], // :path
				headers[3], // :authority
			}

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers[1:]),
			}

			conn.WriteHeaders(hp)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	// All HTTP/2 requests MUST include exactly one valid value for
	// the ":method", ":scheme", and ":path" pseudo-header fields,
	// unless it is a CONNECT request (Section 8.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that omits \":scheme\" pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = []hpack.HeaderField{
				headers[0], // :method
				headers[2], // :path
				headers[3], // :authority
			}

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

	// All HTTP/2 requests MUST include exactly one valid value for
	// the ":method", ":scheme", and ":path" pseudo-header fields,
	// unless it is a CONNECT request (Section 8.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that omits \":path\" pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = []hpack.HeaderField{
				headers[0], // :method
				headers[1], // :scheme
				headers[3], // :authority
			}

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

	// All HTTP/2 requests MUST include exactly one valid value for
	// the ":method", ":scheme", and ":path" pseudo-header fields,
	// unless it is a CONNECT request (Section 8.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with duplicated \":method\" pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField(":method", headers[0].Value))

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

	// All HTTP/2 requests MUST include exactly one valid value for
	// the ":method", ":scheme", and ":path" pseudo-header fields,
	// unless it is a CONNECT request (Section 8.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with duplicated \":scheme\" pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField(":scheme", headers[1].Value))

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

	// All HTTP/2 requests MUST include exactly one valid value for
	// the ":method", ":scheme", and ":path" pseudo-header fields,
	// unless it is a CONNECT request (Section 8.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with duplicated \":path\" pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField(":path", headers[2].Value))

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
