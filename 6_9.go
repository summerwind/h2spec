package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"io"
)

func WindowUpdateTestGroup() *TestGroup {
	tg := NewTestGroup("6.9", "WINDOW_UPDATE")

	tg.AddTestCase(NewTestCase(
		"Sends a WINDOW_UPDATE frame with an flow control window increment of 0",
		"the endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteWindowUpdate(0, 0)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a WINDOW_UPDATE frame with an flow control window increment of 0 on a stream",
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
			}

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)
			http2Conn.fr.WriteWindowUpdate(1, 0)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a WINDOW_UPDATE frame with a length other than a multiple of 4 octets",
		"the endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x03\x08\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x01")

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestGroup(TheFlowControlWindowTestGroup())
	tg.AddTestGroup(InitialFlowControlWindowSizeTestGroup())

	return tg
}

func TheFlowControlWindowTestGroup() *TestGroup {
	tg := NewTestGroup("6.9.1", "The Flow Control Window")

	tg.AddTestCase(NewTestCase(
		"Sends multiple WINDOW_UPDATE frames on a connection increasing the flow control window to above 2^31-1",
		"the endpoint MUST sends a GOAWAY frame with a FLOW_CONTROL_ERROR code.",
		func(ctx *Context) (expected []Result, actual Result) {
			expected = []Result{
				&ResultFrame{http2.FrameGoAway, FlagDefault, http2.ErrCodeFlowControl},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteWindowUpdate(0, 2147483647)
			http2Conn.fr.WriteWindowUpdate(0, 2147483647)

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					if err == io.EOF {
						actual = &ResultConnectionClose{}
					} else if err == TIMEOUT {
						if actual == nil {
							actual = &ResultTestTimeout{}
						}
					} else {
						actual = &ResultError{err}
					}
					break loop
				}

				switch f := f.(type) {
				case *http2.GoAwayFrame:
					actual = &ResultFrame{f.Header().Type, FlagDefault, f.ErrCode}
					if f.ErrCode == http2.ErrCodeFlowControl {
						break loop
					}
				default:
					actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
				}
			}

			return expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends multiple WINDOW_UPDATE frames on a stream increasing the flow control window to above 2^31-1",
		"the endpoint MUST sends a RST_STREAM with the error code of FLOW_CONTROL_ERROR code.",
		func(ctx *Context) (expected []Result, actual Result) {
			expected = []Result{
				&ResultFrame{http2.FrameRSTStream, FlagDefault, http2.ErrCodeFlowControl},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
			}

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteWindowUpdate(1, 2147483647)
			http2Conn.fr.WriteWindowUpdate(1, 2147483647)

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					if err == io.EOF {
						actual = &ResultConnectionClose{}
					} else if err == TIMEOUT {
						if actual == nil {
							actual = &ResultTestTimeout{}
						}
					} else {
						actual = &ResultError{err}
					}
					break loop
				}

				switch f := f.(type) {
				case *http2.RSTStreamFrame:
					actual = &ResultFrame{f.Header().Type, FlagDefault, f.ErrCode}
					if f.ErrCode == http2.ErrCodeFlowControl {
						break loop
					}
				default:
					actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
				}
			}

			return expected, actual
		},
	))

	return tg
}

func InitialFlowControlWindowSizeTestGroup() *TestGroup {
	tg := NewTestGroup("6.9.2", "Initial Flow Control Window Size")

	tg.AddTestCase(NewTestCase(
		"Sends a SETTINGS_INITIAL_WINDOW_SIZE settings with an exceeded maximum window size value",
		"the endpoint MUST respond with a connection error of type FLOW_CONTROL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x04\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x00\x04\x80\x00\x00\x00")

			actualCodes := []http2.ErrCode{http2.ErrCodeFlowControl}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
