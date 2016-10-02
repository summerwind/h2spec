package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StreamDependencies() *spec.TestGroup {
	tg := NewTestGroup("5.3.1", "Stream Dependencies")

	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends HEADERS frame that depend on itself",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      2,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
				Priority: http2.PriorityParam{
					StreamDep: 3,
					Exclusive: false,
					Weight:    255,
				},
			}
			conn.WriteHeaders(hp)

			return spec.VerifyStreamError(conn, http2.ErrCodeProtocol)
		},
	})

	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends PRIORITY frame that depend on itself",
		Requirement: "The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			var streamID uint32 = 2
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
