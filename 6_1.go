package h2spec

import (
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func DataTestGroup() *TestGroup {
	tg := NewTestGroup("6.1", "DATA")

	tg.AddTestCase(NewTestCase(
		"Sends a DATA frame with 0x0 stream identifier",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteData(0, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a DATA frame on the stream that is not opend",
		"The endpoint MUST respond with a stream error of type STREAM_CLOSED.",
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
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)
			http2Conn.fr.WriteData(1, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeStreamClosed}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	/*
		tg.AddTestCase(NewTestCase(
			"Sends a DATA frame with invalid pad length",
			"The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR."
			func(ctx *Context) (expected []Result, actual Result) {},
		))
	*/

	return tg
}
