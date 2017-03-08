package client

import (
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Priority() *spec.ClientTestGroup {
	tg := NewTestGroup("6.3", "PRIORITY")

	// The PRIORITY frame always identifies a stream. If a PRIORITY
	// frame is received with a stream identifier of 0x0, the recipient
	// MUST respond with a connection error (Section 5.4.1) of type
	// PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a PRIORITY frame with 0x0 stream identifier",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			pp := http2.PriorityParam{
				StreamDep: 0,
				Exclusive: false,
				Weight:    255,
			}
			conn.WritePriority(0, pp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// A PRIORITY frame with a length other than 5 octets MUST be
	// treated as a stream error (Section 5.4.2) of type
	// FRAME_SIZE_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a PRIORITY frame with a length other than 5 octets",
		Requirement: "The endpoint MUST respond with a stream error of type FRAME_SIZE_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			// PRIORITY frame:
			// length: 4, flags: 0x0, stream_id: 0x01
			var flags http2.Flags
			conn.WriteRawFrame(http2.FramePriority, flags, req.StreamID, []byte("\x80\x00\x00\x01"))

			return spec.VerifyStreamError(conn, http2.ErrCodeFrameSize)
		},
	})

	return tg
}
