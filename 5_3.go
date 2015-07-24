package h2spec

import (
	"github.com/bradfitz/http2"
)

func StreamPriorityTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("5.3", "Stream Priority")

	tg.AddTestGroup(StreamDependenciesTestGroup(ctx))

	return tg
}

func StreamDependenciesTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("5.3.1", "Stream Dependencies")

	tg.AddTestCase(NewTestCase(
		"Sends HEADERS frame that depend on itself",
		"The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			var pp http2.PriorityParam
			pp.StreamDep = 3
			pp.Exclusive = false
			pp.Weight = 255

			var hp http2.HeadersFrameParam
			hp.StreamID = 3
			hp.EndStream = true
			hp.EndHeaders = true
			hp.Priority = pp
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends PRIORITY frame that depend on itself",
		"The endpoint MUST treat this as a stream error of type PROTOCOL_ERROR",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			var pp http2.PriorityParam
			pp.StreamDep = 2
			pp.Exclusive = false
			pp.Weight = 255

			http2Conn.fr.WritePriority(2, pp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
