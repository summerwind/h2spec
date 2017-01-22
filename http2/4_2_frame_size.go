package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func FrameSize() *spec.TestGroup {
	tg := NewTestGroup("4.2", "Frame Size")

	// All implementations MUST be capable of receiving and minimally
	// processing frames up to 2^14 octets in length, plus the 9-octet
	// frame header (Section 4.1).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a DATA frame with 2^14 octets in length",
		Requirement: "The endpoint MUST be capable of receiving and minimally processing frames up to 2^14 octets in length.",
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

			data := spec.DummyString(conn.MaxFrameSize())
			conn.WriteData(streamID, true, []byte(data))

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// An endpoint MUST send an error code of FRAME_SIZE_ERROR
	// if a frame exceeds the size defined in SETTINGS_MAX_FRAME_SIZE,
	// exceeds any limit defined for the frame type, or is too small
	// to contain mandatory frame data.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a large size DATA frame that exceeds the SETTINGS_MAX_FRAME_SIZE",
		Requirement: "The endpoint MUST send an error code of FRAME_SIZE_ERROR.",
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

			data := spec.DummyString(conn.MaxFrameSize() + 1)
			conn.WriteData(streamID, true, []byte(data))

			return spec.VerifyStreamError(conn, http2.ErrCodeFrameSize)
		},
	})

	// A frame size error in a frame that could alter the state of
	// the entire connection MUST be treated as a connection error
	// (Section 5.4.1); this includes any frame carrying a header block
	// (Section 4.3) (that is, HEADERS, PUSH_PROMISE, and CONTINUATION),
	// SETTINGS, and any frame with a stream identifier of 0.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a large size HEADERS frame that exceeds the SETTINGS_MAX_FRAME_SIZE",
		Requirement: "The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			headers = append(headers, spec.DummyHeaders(c, 5)...)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeFrameSize)
		},
	})

	return tg
}
