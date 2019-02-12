package generic

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StreamsAndMultiplexing() *spec.TestGroup {
	tg := NewTestGroup("2", "Streams and Multiplexing")

	// RFC7540, 5.1, idle:
	// Receiving any frame other than HEADERS or PRIORITY on a stream
	// in this state MUST be treated as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a PRIORITY frame on idle stream",
		Requirement: "The endpoint MUST accept PRIORITY frame.",
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
			conn.WritePriority(1, pp)

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	// RFC7540, 5.1, half-closed (remote):
	// If an endpoint receives additional frames, other than
	// WINDOW_UPDATE, PRIORITY, or RST_STREAM, for a stream that is in
	// this state, it MUST respond with a stream error (Section 5.4.2)
	// of type STREAM_CLOSED.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a WINDOW_UPDATE frame on half-closed (remote) stream",
		Requirement: "The endpoint MUST accept WINDOW_UPDATE frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Set INITIAL_WINDOW_SIZE to zero to prevent the peer from
			// closing the stream.
			settings := http2.Setting{
				ID:  http2.SettingInitialWindowSize,
				Val: 0,
			}
			conn.WriteSettings(settings)

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			conn.WriteWindowUpdate(streamID, 1)

			return spec.VerifyEventType(conn, spec.EventDataFrame)
		},
	})

	// RFC7540, 5.1, half-closed (remote):
	// If an endpoint receives additional frames, other than
	// WINDOW_UPDATE, PRIORITY, or RST_STREAM, for a stream that is in
	// this state, it MUST respond with a stream error (Section 5.4.2)
	// of type STREAM_CLOSED.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a PRIORITY frame on half-closed (remote) stream",
		Requirement: "The endpoint MUST accept PRIORITY frame.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Set INITIAL_WINDOW_SIZE to zero to prevent the peer from
			// closing the stream.
			settings := http2.Setting{
				ID:  http2.SettingInitialWindowSize,
				Val: 0,
			}
			conn.WriteSettings(settings)

			headers := spec.CommonHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			pp := http2.PriorityParam{
				StreamDep: 0,
				Exclusive: false,
				Weight:    255,
			}
			conn.WritePriority(streamID, pp)

			conn.WriteWindowUpdate(streamID, 1)

			return spec.VerifyEventType(conn, spec.EventDataFrame)
		},
	})

	// RFC7540, 5.1, half-closed (remote):
	// If an endpoint receives additional frames, other than
	// WINDOW_UPDATE, PRIORITY, or RST_STREAM, for a stream that is in
	// this state, it MUST respond with a stream error (Section 5.4.2)
	// of type STREAM_CLOSED.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a RST_STREAM frame on half-closed (remote) stream",
		Requirement: "The endpoint MUST accept RST_STREAM frame.",
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

			conn.WriteRSTStream(streamID, http2.ErrCodeCancel)

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	// RFC7540, 5.1, closed:
	// An endpoint MUST NOT send frames other than PRIORITY on a closed
	// stream. An endpoint that receives any frame other than PRIORITY
	// after receiving a RST_STREAM MUST treat that as a stream error
	// (Section 5.4.2) of type STREAM_CLOSED.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a PRIORITY frame on closed stream",
		Requirement: "The endpoint MUST accept PRIORITY frame.",
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

			err = spec.VerifyStreamClose(conn)
			if err != nil {
				return err
			}

			pp := http2.PriorityParam{
				StreamDep: 0,
				Exclusive: false,
				Weight:    255,
			}
			conn.WritePriority(streamID, pp)

			data := [8]byte{}
			conn.WritePing(false, data)

			return spec.VerifyPingFrameWithAck(conn, data)
		},
	})

	return tg
}
