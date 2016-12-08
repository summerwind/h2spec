package h2spec

import (
	"fmt"
	"io"
	"net"
	"syscall"

	"golang.org/x/net/http2"
)

func PriorityTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.3", "PRIORITY")

	tg.AddTestCase(NewTestCase(
		"Sends a PRIORITY frame with 0x0 stream identifier",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
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

			// PRIORITY Frame
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x05\x02\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x80\x00\x00\x01\x0a")

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a PRIORITY frame with a length other than 5 octets",
		"The endpoint MUST respond with a stream error of type FRAME_SIZE_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			// Set INITIAL_WINDOW_SIZE to zero to prevent the peer from closing the stream
			settings := http2.Setting{http2.SettingInitialWindowSize, 0}
			http2Conn.fr.WriteSettings(settings)

			hdrs := commonHeaderFields(ctx)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			// PRIORITY Frame
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x04\x02\x00\x00\x00\x00\x01")
			fmt.Fprintf(http2Conn.conn, "\x80\x00\x00\x01")

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a PRIORITY frame for an idle stream, then send a HEADER frame for a lower stream id",
		"The endpoint MUST respond to the HEADER request with no connection error.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			// PRIORITY Frame
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x05\x02\x00\x00\x00\x00\x03")
			fmt.Fprintf(http2Conn.conn, "\x80\x00\x00\x00\x0a")

			hdrs2 := commonHeaderFields(ctx)
			hdrs2[0].Value = "HEAD"

			var hp2 http2.HeadersFrameParam
			hp2.StreamID = 1
			hp2.EndStream = true
			hp2.EndHeaders = true
			hp2.BlockFragment = http2Conn.EncodeHeader(hdrs2)
			http2Conn.fr.WriteHeaders(hp2)

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
					actual = CreateResultFrame(f)
					if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) && f.Header().Length == 0 {
						pass = true
						break loop
					}
				case *http2.HeadersFrame:
					actual = CreateResultFrame(f)
					if f.FrameHeader.Flags.Has(http2.FlagHeadersEndStream) {
						pass = true
						break loop
					}
				case *http2.GoAwayFrame:
					actual = CreateResultFrame(f)
					break loop
				case *http2.RSTStreamFrame:
					actual = CreateResultFrame(f)
					break loop
				}
				actual = CreateResultFrame(f)
			}

			return pass, expected, actual
		},
	))

	return tg
}
