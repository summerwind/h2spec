package h2spec

import (
	"golang.org/x/net/http2"
)

func FrameSizeTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("4.2", "Frame Size")

	tg.AddTestCase(NewTestCase(
		"Sends large size frame that exceeds the SETTINGS_MAX_FRAME_SIZE",
		"The endpoint MUST send a FRAME_SIZE_ERROR error.",
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
			max_size, ok := http2Conn.Settings[http2.SettingMaxFrameSize]
			if !ok {
				max_size = 18384
			}

			http2Conn.fr.WriteData(1, true, []byte(dummyData(int(max_size)+1)))

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
