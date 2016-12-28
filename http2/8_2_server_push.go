package http2

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func ServerPush() *spec.TestGroup {
	tg := NewTestGroup("8.2", "Server Push")

	// A client cannot push. Thus, servers MUST treat the receipt of
	// a PUSH_PROMISE frame as a connection error (Section 5.4.1) of
	// type PROTOCOL_ERROR. Clients MUST reject any attempt to change
	// the SETTINGS_ENABLE_PUSH setting to a value other than 0 by
	// treating the message as a connection error (Section 5.4.1) of
	// type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a PUSH_PROMISE frame",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonHeaders(c)

			pp := http2.PushPromiseParam{
				StreamID:      streamID,
				PromiseID:     3,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}

			conn.WritePushPromise(pp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	return tg
}
