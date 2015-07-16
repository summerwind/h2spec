package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
)

func GoawayTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.8", "GOAWAY")

	tg.AddTestCase(NewTestCase(
		"Sends a GOAWAY frame with the stream identifier that is not 0x0",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x08\x07\x00\x00\x00\x00\x03")
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x00\x00\x00\x00\x00")

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
