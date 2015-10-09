package h2spec

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"io"
	"net"
	"syscall"
)

func HttpRequestResponseExchangeTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("8.1", "HTTP Request/Response Exchange")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame as HEAD request",
		"The endpoint should respond with no DATA frame or empty DATA frame.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultFrame{LengthDefault, http2.FrameHeaders, http2.FlagHeadersEndStream, ErrCodeDefault},
				&ResultFrame{0, http2.FrameData, http2.FlagDataEndStream, ErrCodeDefault},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs[0].Value = "HEAD"

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						rf, ok := actual.(*ResultFrame)
						if actual == nil || (ok && rf.Type != http2.FrameGoAway) {
							actual = &ResultConnectionClose{}
						}
					} else if err == TIMEOUT {
						if actual == nil {
							actual = &ResultTestTimeout{}
						}
					} else {
						actual = &ResultError{err}
					}
					break loop
				}

				switch f := f.(type) {
				case *http2.DataFrame:
					actual = CreateResultFrame(f)
					if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) && f.Header().Length == 0 {
						pass = true
						break loop
					}
				case *http2.HeadersFrame:
					actual = CreateResultFrame(f)
					if f.FrameHeader.Flags.Has(http2.FlagHeadersEndStream) {
						pass = true
						break loop
					}
				}
				actual = CreateResultFrame(f)
			}

			return pass, expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame containing trailer part",
		"The endpoint should respond with HEADERS frame.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultFrame{LengthDefault, http2.FrameHeaders, http2.FlagHeadersEndStream, ErrCodeDefault},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs[0].Value = "POST"
			hdrs = append(hdrs, pair("content-length", "4"))
			hdrs = append(hdrs, pair("content-type", "text/plain"))
			hdrs = append(hdrs, pair("trailer", "x-test"))

			var hp1 http2.HeadersFrameParam
			hp1.StreamID = 1
			hp1.EndStream = false
			hp1.EndHeaders = true
			hp1.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp1)

			http2Conn.fr.WriteData(1, false, []byte("test"))

			trailers := []hpack.HeaderField{
				pair("x-test", "ok"),
			}

			var hp2 http2.HeadersFrameParam
			hp2.StreamID = 1
			hp2.EndStream = true
			hp2.EndHeaders = true
			hp2.BlockFragment = http2Conn.EncodeHeader(trailers)
			http2Conn.fr.WriteHeaders(hp2)

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						rf, ok := actual.(*ResultFrame)
						if actual == nil || (ok && rf.Type != http2.FrameGoAway) {
							actual = &ResultConnectionClose{}
						}
					} else if err == TIMEOUT {
						if actual == nil {
							actual = &ResultTestTimeout{}
						}
					} else {
						actual = &ResultError{err}
					}
					break loop
				}

				switch f := f.(type) {
				case *http2.HeadersFrame:
					actual = CreateResultFrame(f)
					if f.FrameHeader.Flags.Has(http2.FlagHeadersEndHeaders) {
						pass = true
						break loop
					}
				}
				actual = CreateResultFrame(f)
			}

			return pass, expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a second HEADERS frame without the END_STREAM flag",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs[0].Value = "POST"
			hdrs = append(hdrs, pair("trailer", "x-test"))

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteData(1, false, []byte("test"))

			trailers := []hpack.HeaderField{
				pair("x-test", "ok"),
			}
			hp.BlockFragment = http2Conn.EncodeHeader(trailers)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestGroup(HttpHeaderFieldsTestGroup(ctx))

	return tg
}

func HttpHeaderFieldsTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("8.1.2", "HTTP Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the header field name in uppercase letters",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("X-TEST", "test"))

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestGroup(PseudoHeaderFieldsTestGroup(ctx))
	tg.AddTestGroup(ConnectionSpecificHeaderFieldsTestGroup(ctx))
	tg.AddTestGroup(RequestPseudoHeaderFieldsTestGroup(ctx))
	tg.AddTestGroup(MalformedRequestsAndResponsesTestGroup(ctx))

	return tg
}

func PseudoHeaderFieldsTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("8.1.2.1", "Pseudo-Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the pseudo-header field defined for response",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair(":status", "200"))

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the invalid pseudo-header field",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair(":test", "test"))

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains a pseudo-header field that appears in a header block after a regular header field",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			tmp := []hpack.HeaderField{
				pair("x-test", "test"),
			}
			hdrs = append(tmp, hdrs...)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}

func ConnectionSpecificHeaderFieldsTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("8.1.2.2", "Connection-Specific Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the connection-specific header field",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("connection", "keep-alive"))

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the TE header field that contain any value other than \"trailers\"",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("trailers", "test"))
			hdrs = append(hdrs, pair("te", "trailers, deflate"))

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}

func RequestPseudoHeaderFieldsTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("8.1.2.3", "Request Pseudo-Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that omits mandatory pseudo-header fields",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			tmp := hdrs[0:2]
			hdrs = append(tmp, hdrs[3])

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame containing more than one pseudo-header fields with the same name",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs1 := commonHeaderFields(ctx)
			hdrs2 := commonHeaderFields(ctx)
			hdrs := append(hdrs1, hdrs2...)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}

func MalformedRequestsAndResponsesTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("8.1.2.6", "Malformed Requests and Responses")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the \"content-length\" header field which does not equal the sum of the DATA frame payload lengths",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("content-length", "1"))
			hdrs[0].Value = "POST"

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)
			http2Conn.fr.WriteData(1, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the \"content-length\" header field which does not equal the sum of the multiple DATA frame payload lengths",
		"The endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("content-length", "1"))
			hdrs[0].Value = "POST"

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = false
			hp.EndHeaders = true
			hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
			http2Conn.fr.WriteHeaders(hp)
			http2Conn.fr.WriteData(1, false, []byte("test"))
			http2Conn.fr.WriteData(1, true, []byte("test"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestStreamError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
