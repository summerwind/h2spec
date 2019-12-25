package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StreamDependencies() *spec.TestGroup {
	tg := NewTestGroup("5.3.1", "Stream Dependencies")

	// A stream cannot depend on itself. An endpoint MUST treat this
	// as a stream error (Section 5.4.2) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends HEADERS frame that depends on itself",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR.",
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
				Priority: http2.PriorityParam{
					StreamDep: streamID,
					Exclusive: false,
					Weight:    255,
				},
			}
			conn.WriteHeaders(hp)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	// A stream cannot depend on itself. An endpoint MUST treat this
	// as a stream error (Section 5.4.2) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends PRIORITY frame that depend on itself",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			priorityParam := http2.PriorityParam{
				StreamDep: streamID,
				Exclusive: false,
				Weight:    255,
			}
			conn.WritePriority(streamID, priorityParam)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
