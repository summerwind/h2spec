package h2spec

import (
	"bytes"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func ServerPushTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("8.2", "Server Push")

	tg.AddTestCase(NewTestCase(
		"Sends a PUSH_PROMISE frame",
		"The endpoint MUST treat the receipt of a PUSH_PROMISE frame as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			var buf bytes.Buffer
			hdrs := commonHeaderFields(ctx)
			enc := hpack.NewEncoder(&buf)
			for _, hf := range hdrs {
				_ = enc.WriteField(hf)
			}

			var pp http2.PushPromiseParam
			pp.StreamID = 1
			pp.PromiseID = 3
			pp.EndHeaders = true
			pp.BlockFragment = buf.Bytes()
			http2Conn.fr.WritePushPromise(pp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
