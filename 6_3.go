package h2spec

import (
	"fmt"
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

	return tg
}
