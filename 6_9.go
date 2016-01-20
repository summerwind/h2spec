package h2spec

import (
	"errors"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"net"
	"syscall"
)

func WindowUpdateTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.9", "WINDOW_UPDATE")

	tg.AddTestCase(NewTestCase(
		"Sends a WINDOW_UPDATE frame",
		"The endpoint is expected to send the DATA frame based on the window size.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultFrame{LengthDefault, http2.FrameData, FlagDefault, ErrCodeDefault},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			settings := http2.Setting{http2.SettingInitialWindowSize, 1}
			http2Conn.fr.WriteSettings(settings)

			hdrs := commonHeaderFields(ctx)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			winUpdated := false

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						rf, ok := actual.(*ResultFrame)
						if actual == nil || (ok && rf.Type != http2.FrameGoAway) {
							actual = &ResultConnectionClose{}
						}
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
				case *http2.DataFrame:
					if winUpdated {
						// Let's skip this test if the DATA frame has END_STREAM flag.
						if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) {
							actual = &ResultSkipped{"The length of DATA frame is 0."}
							return true, nil, actual
						}

						if f.FrameHeader.Length != 1 {
							err := errors.New("The length of DATA frame is invalid.")
							actual = &ResultError{err}
							break loop
						}

						actual = CreateResultFrame(f)
						pass = true
						break loop
					} else {
						// Let's skip this test if the DATA frame has END_STREAM flag.
						if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) {
							actual = &ResultSkipped{"The length of DATA frame is 0."}
							return true, nil, actual
						}

						if f.FrameHeader.Length != 1 {
							err := errors.New("The length of DATA frame is invalid.")
							actual = &ResultError{err}
							break loop
						}

						http2Conn.fr.WriteWindowUpdate(1, 1)
						winUpdated = true
					}
				case *http2.GoAwayFrame:
					actual = CreateResultFrame(f)
					break loop
				case *http2.RSTStreamFrame:
					actual = CreateResultFrame(f)
					break loop
				default:
					actual = CreateResultFrame(f)
				}
			}

			return pass, expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a WINDOW_UPDATE frame with an flow control window increment of 0",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteWindowUpdate(0, 0)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a WINDOW_UPDATE frame with an flow control window increment of 0 on a stream",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

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
		"The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x03\x08\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x01")

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestGroup(TheFlowControlWindowTestGroup(ctx))
	tg.AddTestGroup(InitialFlowControlWindowSizeTestGroup(ctx))

	return tg
}

func TheFlowControlWindowTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.9.1", "The Flow Control Window")

	tg.AddTestCase(NewTestCase(
		"Sends multiple WINDOW_UPDATE frames on a connection increasing the flow control window to above 2^31-1",
		"The endpoint MUST sends a GOAWAY frame with a FLOW_CONTROL_ERROR code.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultFrame{LengthDefault, http2.FrameGoAway, FlagDefault, http2.ErrCodeFlowControl},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteWindowUpdate(0, 2147483647)
			http2Conn.fr.WriteWindowUpdate(0, 2147483647)

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						rf, ok := actual.(*ResultFrame)
						if actual == nil || (ok && rf.Type != http2.FrameGoAway) {
							actual = &ResultConnectionClose{}
						}
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
					actual = CreateResultFrame(f)
					if f.ErrCode == http2.ErrCodeFlowControl {
						pass = true
					}
					break loop
				default:
					actual = CreateResultFrame(f)
				}
			}

			return pass, expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends multiple WINDOW_UPDATE frames on a stream increasing the flow control window to above 2^31-1",
		"The endpoint MUST send a RST_STREAM with the error code of FLOW_CONTROL_ERROR code.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultFrame{LengthDefault, http2.FrameRSTStream, FlagDefault, http2.ErrCodeFlowControl},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

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
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						rf, ok := actual.(*ResultFrame)
						if actual == nil || (ok && rf.Type != http2.FrameGoAway) {
							actual = &ResultConnectionClose{}
						}
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
					actual = CreateResultFrame(f)
					break loop
				case *http2.RSTStreamFrame:
					actual = CreateResultFrame(f)
					if f.ErrCode == http2.ErrCodeFlowControl {
						pass = true
					}
					break loop
				default:
					actual = CreateResultFrame(f)
				}
			}

			return pass, expected, actual
		},
	))

	return tg
}

func InitialFlowControlWindowSizeTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.9.2", "Initial Flow Control Window Size")

	tg.AddTestCase(NewTestCase(
		"Sends a SETTINGS_INITIAL_WINDOW_SIZE settings with an exceeded maximum window size value",
		"The endpoint MUST respond with a connection error of type FLOW_CONTROL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
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
