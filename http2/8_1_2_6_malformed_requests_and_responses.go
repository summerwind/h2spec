package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func MalformedRequestsAndResponses() *spec.TestGroup {
	tg := NewTestGroup("8.1.2.6", "Malformed Requests and Responses")

	// A request or response that includes a payload body can include
	// a content-length header field. A request or response is also
	// malformed if the value of a content-length header field does
	// not equal the sum of the DATA frame payload lengths that form
	// the body. A response that is defined to have no payload, as
	// described in [RFC7230], Section 3.3.2, can have a non-zero
	// content-length header field, even though no content is included
	// in DATA frames.
	//
	// Intermediaries that process HTTP requests or responses (i.e.,
	// any intermediary not acting as a tunnel) MUST NOT forward a
	// malformed request or response. Malformed requests or responses
	// that are detected MUST be treated as a stream error
	// (Section 5.4.2) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with the \"content-length\" header field which does not equal the DATA frame payload length",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"
			headers = append(headers, spec.HeaderField("content-length", "1"))

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)
			conn.WriteData(streamID, true, []byte("test"))

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	// A request or response that includes a payload body can include
	// a content-length header field. A request or response is also
	// malformed if the value of a content-length header field does
	// not equal the sum of the DATA frame payload lengths that form
	// the body. A response that is defined to have no payload, as
	// described in [RFC7230], Section 3.3.2, can have a non-zero
	// content-length header field, even though no content is included
	// in DATA frames.
	//
	// Intermediaries that process HTTP requests or responses (i.e.,
	// any intermediary not acting as a tunnel) MUST NOT forward a
	// malformed request or response. Malformed requests or responses
	// that are detected MUST be treated as a stream error
	// (Section 5.4.2) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame with the \"content-length\" header field which does not equal the sum of the multiple DATA frames payload length",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"
			headers = append(headers, spec.HeaderField("content-length", "1"))

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)
			conn.WriteData(streamID, false, []byte("test"))
			conn.WriteData(streamID, true, []byte("test"))

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
