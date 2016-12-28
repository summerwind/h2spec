package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func PseudoHeaderFields() *spec.TestGroup {
	tg := NewTestGroup("8.1.2.1", "Pseudo-Header Fields")

	// Pseudo-header fields are only valid in the context in which
	// they are defined. Pseudo-header fields defined for requests
	// MUST NOT appear in responses; pseudo-header fields defined
	// for responses MUST NOT appear in requests. Pseudo-header
	// fields MUST NOT appear in trailers. Endpoints MUST treat
	// a request or response that contains undefined or invalid
	// pseudo-header fields as malformed (Section 8.1.2.6).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that contains a unknown pseudo-header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField(":test", "ok"))

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

	// Pseudo-header fields are only valid in the context in which
	// they are defined. Pseudo-header fields defined for requests
	// MUST NOT appear in responses; pseudo-header fields defined
	// for responses MUST NOT appear in requests. Pseudo-header
	// fields MUST NOT appear in trailers. Endpoints MUST treat
	// a request or response that contains undefined or invalid
	// pseudo-header fields as malformed (Section 8.1.2.6).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that contains the pseudo-header field defined for response",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.HeaderField(":status", "200"))

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

	// Pseudo-header fields are only valid in the context in which
	// they are defined. Pseudo-header fields defined for requests
	// MUST NOT appear in responses; pseudo-header fields defined
	// for responses MUST NOT appear in requests. Pseudo-header
	// fields MUST NOT appear in trailers. Endpoints MUST treat
	// a request or response that contains undefined or invalid
	// pseudo-header fields as malformed (Section 8.1.2.6).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that contains a pseudo-header field as trailers",
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
				spec.HeaderField(":method", "POST"),
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

	// All pseudo-header fields MUST appear in the header block before
	// regular header fields. Any request or response that contains
	// a pseudo-header field that appears in a header block after
	// a regular header field MUST be treated as malformed
	// (Section 8.1.2.6).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame that contains a pseudo-header field that appears in a header block after a regular header field",
		Requirement: "The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := []hpack.HeaderField{
				spec.HeaderField("x-test", "ok"),
			}
			headers = append(headers, spec.CommonHeaders(c)...)

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
