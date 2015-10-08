package h2spec

import (
	"golang.org/x/net/http2"
)

func StreamStatesTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("5.1", "Stream States")

	tg.AddTestCase(NewTestCase(
		"idle: Sends a DATA frame",
		"The endpoint MUST treat this as a connection error (Section 5.4.1) of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteData(1, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"idle: Sends a RST_STREAM frame",
		"The endpoint MUST treat this as a connection error (Section 5.4.1) of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteRSTStream(1, http2.ErrCodeCancel)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"idle: Sends a WINDOW_UPDATE frame",
		"The endpoint MUST treat this as a connection error (Section 5.4.1) of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			http2Conn.fr.WriteWindowUpdate(1, 100)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"idle: Sends a CONTINUATION frame",
		"The endpoint MUST treat this as a connection error (Section 5.4.1) of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			blockFragment := http2Conn.EncodeHeader(hdrs)

			http2Conn.fr.WriteContinuation(1, true, blockFragment)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"half closed (remote): Sends a DATA frame",
		"The endpoint MUST respond with a stream error (Section 5.4.2) of type STREAM_CLOSED.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = blockFragment
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteData(1, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeStreamClosed}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"half closed (remote): Sends a HEADERS frame",
		"The endpoint MUST respond with a stream error (Section 5.4.2) of type STREAM_CLOSED.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp1 http2.HeadersFrameParam
			hp1.StreamID = 1
			hp1.EndStream = true
			hp1.EndHeaders = true
			hp1.BlockFragment = blockFragment
			http2Conn.fr.WriteHeaders(hp1)

			var hp2 http2.HeadersFrameParam
			hp2.StreamID = 1
			hp2.EndStream = true
			hp2.EndHeaders = true
			hp2.BlockFragment = blockFragment
			http2Conn.fr.WriteHeaders(hp2)

			actualCodes := []http2.ErrCode{http2.ErrCodeStreamClosed}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"half closed (remote): Sends a CONTINUATION frame",
		"The endpoint MUST respond with a stream error (Section 5.4.2) of type STREAM_CLOSED.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = blockFragment
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteContinuation(1, true, blockFragment)

			actualCodes := []http2.ErrCode{http2.ErrCodeStreamClosed, http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	if ctx.Strict {
		tg.AddTestCase(NewTestCase(
			"closed: Sends a DATA frame",
			"The endpoint MUST treat this as a stream error (Section 5.4.2) of type STREAM_CLOSED.",
			func(ctx *Context) (pass bool, expected []Result, actual Result) {
				http2Conn := CreateHttp2Conn(ctx, true)
				defer http2Conn.conn.Close()

				hdrs := commonHeaderFields(ctx)
				blockFragment := http2Conn.EncodeHeader(hdrs)

				var hp http2.HeadersFrameParam
				hp.StreamID = 1
				hp.EndStream = true
				hp.EndHeaders = true
				hp.BlockFragment = blockFragment
				http2Conn.fr.WriteHeaders(hp)

				pass, expected, actual = TestStreamClose(ctx, http2Conn)
				if !pass {
					return pass, expected, actual
				}

				http2Conn.fr.WriteData(1, true, []byte("test"))

				actualCodes := []http2.ErrCode{http2.ErrCodeStreamClosed}
				return TestStreamError(ctx, http2Conn, actualCodes)
			},
		))
	}

	if ctx.Strict {
		tg.AddTestCase(NewTestCase(
			"closed: Sends a HEADERS frame",
			"The endpoint MUST treat this as a stream error (Section 5.4.2) of type STREAM_CLOSED.",
			func(ctx *Context) (pass bool, expected []Result, actual Result) {
				http2Conn := CreateHttp2Conn(ctx, true)
				defer http2Conn.conn.Close()

				hdrs := commonHeaderFields(ctx)
				blockFragment := http2Conn.EncodeHeader(hdrs)

				var hp http2.HeadersFrameParam
				hp.StreamID = 1
				hp.EndStream = true
				hp.EndHeaders = true
				hp.BlockFragment = blockFragment
				http2Conn.fr.WriteHeaders(hp)

				pass, expected, actual = TestStreamClose(ctx, http2Conn)
				if !pass {
					return pass, expected, actual
				}

				http2Conn.fr.WriteHeaders(hp)

				actualCodes := []http2.ErrCode{http2.ErrCodeStreamClosed}
				return TestStreamError(ctx, http2Conn, actualCodes)
			},
		))
	}

	tg.AddTestCase(NewTestCase(
		"closed: Sends a CONTINUATION frame",
		"The endpoint MUST treat this as a stream error (Section 5.4.2) of type STREAM_CLOSED.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("x-dummy1", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy2", dummyData(10000)))
			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = false
			hp.BlockFragment = blockFragment[0:16384]
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteContinuation(1, true, blockFragment[16384:])

			pass, expected, actual = TestStreamClose(ctx, http2Conn)
			if !pass {
				return pass, expected, actual
			}

			http2Conn.fr.WriteContinuation(1, true, blockFragment[16384:])

			actualCodes := []http2.ErrCode{http2.ErrCodeStreamClosed, http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestGroup(StreamIdentifiersTestGroup(ctx))
	tg.AddTestGroup(StreamConcurrencyTestGroup(ctx))

	return tg
}

func StreamIdentifiersTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("5.1.1", "Stream Identifiers")

	tg.AddTestCase(NewTestCase(
		"Sends even-numbered stream identifier",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

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

	if ctx.Strict {
		tg.AddTestCase(NewTestCase(
			"Sends stream identifier that is numerically smaller than previous",
			"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
			func(ctx *Context) (pass bool, expected []Result, actual Result) {
				http2Conn := CreateHttp2Conn(ctx, true)
				defer http2Conn.conn.Close()

				hdrs := commonHeaderFields(ctx)

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
	}

	return tg
}

func StreamConcurrencyTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("5.1.2", "Stream Concurrency")

	tg.AddTestCase(NewTestCase(
		"Sends HEADERS frames that causes their advertised concurrent stream limit to be exceeded",
		"The endpoint MUST treat this as a stream error (Section 5.4.2) of type PROTOCOL_ERROR or REFUSED_STREAM",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			// Skip this test when SETTINGS_MAX_CONCURRENT_STREAMS is unlimited.
			_, ok := http2Conn.Settings[http2.SettingMaxConcurrentStreams]
			if !ok {
				actual = &ResultSkipped{"SETTINGS_MAX_CONCURRENT_STREAMS is unlimited."}
				return true, nil, actual
			}

			// Set INITIAL_WINDOW_SIZE to zero to prevent the peer from closing the stream
			settings := http2.Setting{http2.SettingInitialWindowSize, 0}
			http2Conn.fr.WriteSettings(settings)

			hdrs := commonHeaderFields(ctx)
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
