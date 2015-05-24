package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
)

func RstStreamTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.4", "RST_STREAM")

	tg.AddTestCase(NewTestCase(
		"Sends a RST_STREAM frame with 0x0 stream identifier",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteRSTStream(0, http2.ErrCodeCancel)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a RST_STREAM frame on a idle stream",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteRSTStream(1, http2.ErrCodeCancel)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a RST_STREAM frame with a length other than 4 octets",
		"The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			// RST_STREAM Frame
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x03\x03\x00\x00\x00\x00\x01")
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x00")

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
