package h2spec

import (
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func StreamStatesTestGroup() *TestGroup {
	tg := NewTestGroup("5.1", "Stream States")

	tg.AddTestGroup(StreamIdentifiersTestGroup())
	tg.AddTestGroup(StreamConcurrencyTestGroup())

	return tg
}

func StreamIdentifiersTestGroup() *TestGroup {
	tg := NewTestGroup("5.1.1", "Stream Identifiers")

	tg.AddTestCase(NewTestCase(
		"Sends even-numbered stream identifier",
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
			hp.StreamID = 2
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends stream identifier that is numerically smaller than previous",
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

			var hp1 http2.HeadersFrameParam
			hp1.StreamID = 5
			hp1.EndStream = true
			hp1.EndHeaders = true
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

	return tg
}

func StreamConcurrencyTestGroup() *TestGroup {
	tg := NewTestGroup("5.1.2", "Stream Concurrency")

	tg.AddTestCase(NewTestCase(
		"Sends HEADERS frames that causes their advertised concurrent stream limit to be exceeded",
		"The endpoint MUST treat this as a stream error (Section 5.4.2) of type PROTOCOL_ERROR or REFUSED_STREAM",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			// Skip this test when SETTINGS_MAX_CONCURRENT_STREAMS is unlimited.
			_, ok := http2Conn.Settings[http2.SettingMaxConcurrentStreams]
			if !ok {
				actual = &ResultSkipped{"SETTINGS_MAX_CONCURRENT_STREAMS is unlimited."}
				return nil, actual
			}

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
			}
			hbf := http2Conn.EncodeHeader(hdrs)

			var streamID uint32 = 1
			for i := 0; i <= int(http2Conn.Settings[http2.SettingMaxConcurrentStreams]); i++ {
				var hp http2.HeadersFrameParam
				hp.StreamID = streamID
				hp.EndStream = true
				hp.EndHeaders = true
				hp.BlockFragment = hbf
				http2Conn.fr.WriteHeaders(hp)
				streamID += 2
			}

			actualCodes := []http2.ErrCode{
				http2.ErrCodeProtocol,
				http2.ErrCodeRefusedStream,
			}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
