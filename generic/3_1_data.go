package generic

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Data() *spec.TestGroup {
	tg := NewTestGroup("3.1", "DATA")

	// RFC7540, 6.1:
	// DATA frames (type=0x0) convey arbitrary, variable-length
	// sequences of octets associated with a stream. One or more
	// DATA frames are used, for instance, to carry HTTP request
	// or response payloads.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a DATA frame",
		Requirement: "The endpoint MUST accept DATA frame.",
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

	// RFC7540, 6.1:
	// DATA frames (type=0x0) convey arbitrary, variable-length
	// sequences of octets associated with a stream. One or more
	// DATA frames are used, for instance, to carry HTTP request
	// or response payloads.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends multiple DATA frames",
		Requirement: "The endpoint MUST accept multiple DATA frames.",
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
			conn.WriteData(streamID, false, []byte("test"))
			conn.WriteData(streamID, true, []byte("test"))

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC7540, 6.1:
	// DATA frames MAY also contain padding. Padding can be added to
	// DATA frames to obscure the size of messages. Padding is a
	// security feature; see Section 10.7.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a DATA frame with padding",
		Requirement: "The endpoint MUST accept DATA frame with padding.",
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
			conn.WriteDataPadded(streamID, true, []byte("test"), []byte("\x00\x00\x00\x00\x00"))

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	return tg
}
