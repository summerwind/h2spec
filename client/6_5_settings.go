package client

import (
	"errors"
	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func Settings() *spec.ClientTestGroup {
	tg := NewTestGroup("6.5", "SETTINGS")

	// ACK (0x1):
	// When set, bit 0 indicates that this frame acknowledges receipt
	// and application of the peer's SETTINGS frame. When this bit is
	// set, the payload of the SETTINGS frame MUST be empty. Receipt of
	// a SETTINGS frame with the ACK flag set and a length field value
	// other than 0 MUST be treated as a connection error (Section 5.4.1)
	// of type FRAME_SIZE_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a SETTINGS frame with ACK flag and payload",
		Requirement: "The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			_, err := conn.ReadClientPreface()
			if err != nil {
				return err
			}

			_, ok := conn.WaitEventByType(spec.EventSettingsFrame)
			if !ok {
				return errors.New("First frame from client must be SETTINGS")
			}

			// SETTINGS frame:
			// length: 0, flags: 0x1, stream_id: 0x0
			conn.Send([]byte("\x00\x00\x01\x04\x01\x00\x00\x00\x00"))
			conn.Send([]byte("\x00"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeFrameSize)
		},
	})

	// SETTINGS frames always apply to a connection, never a single
	// stream. The stream identifier for a SETTINGS frame MUST be
	// zero (0x0). If an endpoint receives a SETTINGS frame whose
	// stream identifier field is anything other than 0x0, the
	// endpoint MUST respond with a connection error (Section 5.4.1)
	// of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a SETTINGS frame with a stream identifier other than 0x0",
		Requirement: "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			_, err := conn.ReadClientPreface()
			if err != nil {
				return err
			}

			_, ok := conn.WaitEventByType(spec.EventSettingsFrame)
			if !ok {
				return errors.New("First frame from client must be SETTINGS")
			}

			// SETTINGS frame:
			// length: 6, flags: 0x0, stream_id: 0x1
			conn.Send([]byte("\x00\x00\x06\x04\x00\x00\x00\x00\x01"))
			conn.Send([]byte("\x00\x03\x00\x00\x00\x64"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// The SETTINGS frame affects connection state. A badly formed or
	// incomplete SETTINGS frame MUST be treated as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	//
	// A SETTINGS frame with a length other than a multiple of 6 octets
	// MUST be treated as a connection error (Section 5.4.1) of type
	// FRAME_SIZE_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends a SETTINGS frame with a length other than a multiple of 6 octets",
		Requirement: "The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			_, err := conn.ReadClientPreface()
			if err != nil {
				return err
			}

			_, ok := conn.WaitEventByType(spec.EventSettingsFrame)
			if !ok {
				return errors.New("First frame from client must be SETTINGS")
			}

			// SETTINGS frame:
			// length: 3, flags: 0x0, stream_id: 0x0
			conn.Send([]byte("\x00\x00\x03\x04\x00\x00\x00\x00\x00"))
			conn.Send([]byte("\x00\x03\x00"))

			codes := []http2.ErrCode{
				http2.ErrCodeProtocol,
				http2.ErrCodeFrameSize,
			}
			return spec.VerifyStreamError(conn, codes...)
		},
	})

	tg.AddTestGroup(DefinedSETTINGSParameters())
	tg.AddTestGroup(SettingsSynchronization())

	return tg
}
