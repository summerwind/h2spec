package h2spec

import (
	"github.com/bradfitz/http2"
	"io"
	"net"
	"syscall"
)

func ContinuationTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.10", "CONTINUATION")

	tg.AddTestCase(NewTestCase(
		"Sends a CONTINUATION frame",
		"The endpoint must accept the frame.",
		func(ctx *Context) (expected []Result, actual Result) {
			expected = []Result{
				&ResultFrame{http2.FrameHeaders, FlagDefault, ErrCodeDefault},
			}

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

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						actual = &ResultConnectionClose{}
					} else if err == TIMEOUT {
						if actual == nil {
							actual = &ResultTestTimeout{}
						}
					} else {
						actual = &ResultError{err}
					}
					break loop
				}

				actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
				_, ok := f.(*http2.HeadersFrame)
				if ok {
					break loop
				}
			}

			return expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends multiple CONTINUATION frames",
		"The endpoint must accept the frames.",
		func(ctx *Context) (expected []Result, actual Result) {
			expected = []Result{
				&ResultFrame{http2.FrameHeaders, FlagDefault, ErrCodeDefault},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("x-dummy1", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy2", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy3", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy4", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy5", dummyData(10000)))

			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = false
			hp.BlockFragment = blockFragment[0:16384]
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])
			http2Conn.fr.WriteContinuation(1, true, blockFragment[32767:])

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						actual = &ResultConnectionClose{}
					} else if err == TIMEOUT {
						if actual == nil {
							actual = &ResultTestTimeout{}
						}
					} else {
						actual = &ResultError{err}
					}
					break loop
				}

				actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
				_, ok := f.(*http2.HeadersFrame)
				if ok {
					break loop
				}
			}

			return expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a CONTINUATION frame followed by any frame other than CONTINUATION",
		"The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("x-dummy1", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy2", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy3", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy4", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy5", dummyData(10000)))

			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = false
			hp.BlockFragment = blockFragment[0:16384]
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])
			http2Conn.fr.WriteData(1, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a CONTINUATION frame followed by a frame on a different stream",
		"The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("x-dummy1", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy2", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy3", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy4", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy5", dummyData(10000)))

			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = false
			hp.BlockFragment = blockFragment[0:16384]
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])
			http2Conn.fr.WriteContinuation(3, true, blockFragment[32767:])

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a CONTINUATION frame with the stream identifier that is 0x0",
		"The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("x-dummy1", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy2", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy3", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy4", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy5", dummyData(10000)))

			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = false
			hp.BlockFragment = blockFragment[0:16384]
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])
			http2Conn.fr.WriteContinuation(0, true, blockFragment[32767:])

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a CONTINUATION frame after the frame other than HEADERS, PUSH_PROMISE or CONTINUATION",
		"The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteData(1, true, []byte("test"))
			http2Conn.fr.WriteContinuation(1, true, http2Conn.EncodeHeader(hdrs))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
