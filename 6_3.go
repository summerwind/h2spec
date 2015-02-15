package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func PriorityTestGroup() *TestGroup {
	tg := NewTestGroup("6.3", "PRIORITY")

	tg.AddTestCase(NewTestCase(
		"Sends a PRIORITY frame with 0x0 stream identifier",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
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

			// PRIORITY Frame
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x04\x02\x00\x00\x00\x00\x01")
			fmt.Fprintf(http2Conn.conn, "\x80\x00\x00\x01")

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
