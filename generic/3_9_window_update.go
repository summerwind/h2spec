package generic

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func WindowUpdate() *spec.TestGroup {
	tg := NewTestGroup("3.9", "WINDOW_UPDATE")

	// RFC7540, 6.9:
	// The WINDOW_UPDATE frame (type=0x8) is used to implement flow
	// control; see Section 5.2 for an overview.
	//
	// Flow control operates at two levels: on each individual stream
	// and on the entire connection.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a WINDOW_UPDATE frame with stream ID 0",
		Requirement: "The endpoint MUST accept WINDOW_UPDATE frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteWindowUpdate(0, 1)

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	// RFC7540, 6.9:
	// The WINDOW_UPDATE frame (type=0x8) is used to implement flow
	// control; see Section 5.2 for an overview.
	//
	// Flow control operates at two levels: on each individual stream
	// and on the entire connection.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a WINDOW_UPDATE frame with stream ID 1",
		Requirement: "The endpoint MUST accept WINDOW_UPDATE frame.",
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

			conn.WriteWindowUpdate(streamID, 1)

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	return tg
}
