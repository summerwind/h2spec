package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func HeaderCompressionAndDecompression() *spec.TestGroup {
	tg := NewTestGroup("4.3", "Header Compression and Decompression")

	// A decoding error in a header block MUST be treated as
	// a connection error (Section 5.4.1) of type COMPRESSION_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends invalid header block fragment",
		Requirement: "The endpoint MUST terminate the connection with a connection error of type COMPRESSION_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal Header Field with Incremental Indexing without
			// Length and String segment.
			err = conn.Send([]byte("\x00\x00\x01\x01\x05\x00\x00\x00\x01\x40"))
			if err != nil {
				return err
			}

			return spec.VerifyConnectionError(conn, http2.ErrCodeCompression)
		},
	})

	// Each header block is processed as a discrete unit. Header blocks
	// MUST be transmitted as a contiguous sequence of frames, with no
	// interleaved frames of any other type or from any other stream.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a PRIORITY frame while sending the header blocks",
		Requirement: "The endpoint MUST terminate the connection with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			pp := http2.PriorityParam{
				StreamDep: 0,
				Exclusive: false,
				Weight:    255,
			}
			conn.WritePriority(streamID, pp)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// Each header block is processed as a discrete unit. Header blocks
	// MUST be transmitted as a contiguous sequence of frames, with no
	// interleaved frames of any other type or from any other stream.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a HEADERS frame to another stream while sending the header blocks",
		Requirement: "The endpoint MUST terminate the connection with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)

			hp1 := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     false,
				EndHeaders:    false,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp1)

			hp2 := http2.HeadersFrameParam{
				StreamID:      streamID + 2,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp2)

			dummyHeaders := spec.DummyHeaders(c, 1)
			conn.WriteContinuation(streamID, true, conn.EncodeHeaders(dummyHeaders))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
