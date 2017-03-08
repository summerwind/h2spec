package client

import (
	"fmt"

	"golang.org/x/net/http2"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
)

func TheFlowControlWindow() *spec.ClientTestGroup {
	tg := NewTestGroup("6.9.1", "The Flow-Control Window")

	// A sender MUST NOT allow a flow-control window to exceed 2^31-1
	// octets. If a sender receives a WINDOW_UPDATE that causes a
	// flow-control window to exceed this maximum, it MUST terminate
	// either the stream or the connection, as appropriate.
	// For streams, the sender sends a RST_STREAM with an error code
	// of FLOW_CONTROL_ERROR; for the connection, a GOAWAY frame with
	// an error code of FLOW_CONTROL_ERROR is sent.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends multiple WINDOW_UPDATE frames increasing the flow control window to above 2^31-1",
		Requirement: "The endpoint MUST sends a GOAWAY frame with a FLOW_CONTROL_ERROR code.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteWindowUpdate(0, 2147483647)
			conn.WriteWindowUpdate(0, 2147483647)

			actual, passed := conn.WaitEventByType(spec.EventGoAwayFrame)
			switch event := actual.(type) {
			case spec.GoAwayFrameEvent:
				passed = (event.ErrCode == http2.ErrCodeFlowControl)
			default:
				passed = false
			}

			if !passed {
				expected := []string{
					fmt.Sprintf("GOAWAY Frame (Error Code: %s)", http2.ErrCodeFlowControl),
				}

				return &spec.TestError{
					Expected: expected,
					Actual:   actual.String(),
				}
			}

			return nil
		},
	})

	// A sender MUST NOT allow a flow-control window to exceed 2^31-1
	// octets. If a sender receives a WINDOW_UPDATE that causes a
	// flow-control window to exceed this maximum, it MUST terminate
	// either the stream or the connection, as appropriate.
	// For streams, the sender sends a RST_STREAM with an error code
	// of FLOW_CONTROL_ERROR; for the connection, a GOAWAY frame with
	// an error code of FLOW_CONTROL_ERROR is sent.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends multiple WINDOW_UPDATE frames increasing the flow control window to above 2^31-1 on a stream",
		Requirement: "The endpoint MUST sends a RST_STREAM frame with a FLOW_CONTROL_ERROR code.",
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
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			conn.WriteWindowUpdate(req.StreamID, 2147483647)
			conn.WriteWindowUpdate(req.StreamID, 2147483647)

			actual, passed := conn.WaitEventByType(spec.EventRSTStreamFrame)
			switch event := actual.(type) {
			case spec.RSTStreamFrameEvent:
				if event.Header().StreamID == req.StreamID {
					passed = (event.ErrCode == http2.ErrCodeFlowControl)
				}
			default:
				passed = false
			}

			if !passed {
				expected := []string{
					fmt.Sprintf("RST_STREAM Frame (Error Code: %s)", http2.ErrCodeFlowControl),
				}

				return &spec.TestError{
					Expected: expected,
					Actual:   actual.String(),
				}
			}

			return nil
		},
	})

	return tg
}
