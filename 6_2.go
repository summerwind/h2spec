package h2spec

import (
	"bytes"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func HeadersTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.2", "HEADERS")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame followed by any frame other than CONTINUATION",
		"The endpoint MUST treat the receipt of any other type of frame as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = false
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)
			http2Conn.fr.WriteData(1, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame followed by a frame on a different stream",
		"The endpoint MUST treat the receipt of a frame on a different stream as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			var hp1 http2.HeadersFrameParam
			hp1.StreamID = 1
			hp1.EndStream = false
			hp1.EndHeaders = false
			hp1.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp1)

			var hp2 http2.HeadersFrameParam
			hp2.StreamID = 3
			hp2.EndStream = true
			hp2.EndHeaders = true
			hp2.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp2)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame with 0x0 stream identifier",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			var hp http2.HeadersFrameParam
			hp.StreamID = 0
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame with invalid pad length",
		"The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			var buf bytes.Buffer
			hdrs := commonHeaderFields(ctx)
			enc := hpack.NewEncoder(&buf)
			for _, hf := range hdrs {
				_ = enc.WriteField(hf)
			}

			// Payload length: 12, Pad length: 13
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x0c\x01\x0d\x00\x00\x00\x01")
			fmt.Fprintf(http2Conn.conn, "\x0d")
			http2Conn.conn.Write(buf.Bytes())

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
