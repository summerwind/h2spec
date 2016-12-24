package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func FrameFormat() *spec.TestGroup {
	tg := NewTestGroup("4.1", "Frame Format")

	// Type: The 8-bit type of the frame. The frame type determines
	// the format and semantics of the frame. Implementations MUST
	// ignore and discard any frame that has a type that is unknown.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a frame with unknown type",
		Requirement: "The endpoint MUST ignore and discard any frame that has a type that is unknown.",
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

			// UNKONWN Frame:
			// Length: 8, Type: 255, Flags: 0, R: 0, StreamID: 0
			conn.Send("\x00\x00\x08\x16\x00\x00\x00\x00\x00")
			conn.Send("\x00\x00\x00\x00\x00\x00\x00\x00")

			conn.WriteHeaders(hp)

			return spec.VerifyFrameType(conn, http2.FrameHeaders)
		},
	})

	// Flags are assigned semantics specific to the indicated frame
	// type. Flags that have no defined semantics for a particular
	// frame type MUST be ignored and MUST be left unset (0x0) when
	// sending.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a frame with undefined flag",
		Requirement: "The endpoint MUST ignore any flags that is undefined.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// PING Frame:
			// Length: 8, Type: 6, Flags: 255, R: 0, StreamID: 0
			conn.Send("\x00\x00\x08\x06\x16\x00\x00\x00\x00")
			conn.Send("\x00\x00\x00\x00\x00\x00\x00\x00")

			return spec.VerifyFrameType(conn, http2.FramePing)
		},
	})

	// R: A reserved 1-bit field. The semantics of this bit are
	// undefined, and the bit MUST remain unset (0x0) when sending
	// and MUST be ignored when receiving.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a frame with reserved field bit",
		Requirement: "The endpoint MUST ignore the value of reserved field.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			// PING Frame:
			// Length: 8, Type: 6, Flags: 255, R: 1, StreamID: 0
			conn.Send("\x00\x00\x08\x06\x16\x80\x00\x00\x00")
			conn.Send("\x00\x00\x00\x00\x00\x00\x00\x00")

			return spec.VerifyFrameType(conn, http2.FramePing)
		},
	})

	return tg
}
