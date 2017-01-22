package generic

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func HTTPMessageExchanges() *spec.TestGroup {
	tg := NewTestGroup("4", "HTTP Message Exchanges")

	// RFC7540, 8.1:
	// A client sends an HTTP request on a new stream, using a
	// previously unused stream identifier (Section 5.1.1).
	// A server sends an HTTP response on the same stream as the
	// request.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a GET request",
		Requirement: "The endpoint MUST respond to the request.",
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

	// RFC7540, 8.1:
	// A client sends an HTTP request on a new stream, using a
	// previously unused stream identifier (Section 5.1.1).
	// A server sends an HTTP response on the same stream as the
	// request.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEAD request",
		Requirement: "The endpoint MUST respond to the request.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "HEAD"

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

	// RFC7540, 8.1:
	// A client sends an HTTP request on a new stream, using a
	// previously unused stream identifier (Section 5.1.1).
	// A server sends an HTTP response on the same stream as the
	// request.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a POST request",
		Requirement: "The endpoint MUST respond to the request.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)
			conn.WriteData(streamID, true, []byte("test"))

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC7540, 8.1:
	// A client sends an HTTP request on a new stream, using a
	// previously unused stream identifier (Section 5.1.1).
	// A server sends an HTTP response on the same stream as the
	// request.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a POST request with trailers",
		Requirement: "The endpoint MUST respond to the request.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers[0].Value = "POST"
			headers = append(headers, spec.HeaderField("trailer", "x-test"))

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
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(trailers),
			}

			conn.WriteHeaders(hp2)

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	return tg
}
