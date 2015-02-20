package h2spec

import (
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func FrameSizeTestGroup() *TestGroup {
	tg := NewTestGroup("4.2", "Frame Size")

	tg.AddTestCase(NewTestCase(
		"Sends large size frame that exceeds the SETTINGS_MAX_FRAME_SIZE",
		"The endpoint MUST send a FRAME_SIZE_ERROR error.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, false)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteSettings()

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
			http2Conn.fr.WriteData(1, true, []byte(dummyData(16385)))

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
