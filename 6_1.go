package h2spec

import (
	"fmt"
	"golang.org/x/net/http2"
)

func DataTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.1", "DATA")

	tg.AddTestCase(NewTestCase(
		"Sends a DATA frame with 0x0 stream identifier",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteData(0, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a DATA frame on the stream that is not in \"open\" or \"half-closed (local)\" state",
		"The endpoint MUST respond with a stream error (Section 5.4.2) of type STREAM_CLOSED.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs[0].Value = "POST"
			hdrs = append(hdrs, pair("content-length", "4"))

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

	tg.AddTestCase(NewTestCase(
		"Sends a DATA frame with invalid pad length",
		"The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs[0].Value = "POST"
			hdrs = append(hdrs, pair("content-length", "4"))

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			// Data length: 5, Pad length: 6
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x05\x00\x09\x00\x00\x00\x01")
			fmt.Fprintf(http2Conn.conn, "\x06\x54\x65\x73\x74")

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
