package h2spec

import (
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func HttpRequestResponseExchangeTestGroup() *TestGroup {
	tg := NewTestGroup("8.1", "HTTP Request/Response Exchange")

	tg.AddTestGroup(HttpHeaderFieldsTestGroup())

	return tg
}

func HttpHeaderFieldsTestGroup() *TestGroup {
	tg := NewTestGroup("8.1.2", "HTTP Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the header field name in uppercase letters",
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair("X-TEST", "test"),
			}

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

	tg.AddTestGroup(PseudoHeaderFieldsTestGroup())
	tg.AddTestGroup(ConnectionSpecificHeaderFieldsTestGroup())
	tg.AddTestGroup(RequestPseudoHeaderFieldsTestGroup())
	tg.AddTestGroup(MalformedRequestsAndResponsesTestGroup())

	return tg
}

func PseudoHeaderFieldsTestGroup() *TestGroup {
	tg := NewTestGroup("8.1.2.1", "Pseudo-Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the pseudo-header field defined for response",
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair(":status", "200"),
			}

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
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair(":test", "test"),
			}

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
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair("x-test", "test"),
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
			}

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

func ConnectionSpecificHeaderFieldsTestGroup() *TestGroup {
	tg := NewTestGroup("8.1.2.2", "Connection-Specific Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the connection-specific header field",
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair("connection", "keep-alive"),
			}

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
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair("trailers", "test"),
				pair("te", "trailers, deflate"),
			}

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

func RequestPseudoHeaderFieldsTestGroup() *TestGroup {
	tg := NewTestGroup("8.1.2.3", "Request Pseudo-Header Fields")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that is omitted mandatory pseudo-header fields",
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":authority", ctx.Authority()),
			}

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
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "http"),
				pair(":authority", ctx.Authority()),
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "http"),
				pair(":authority", ctx.Authority()),
			}

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

func MalformedRequestsAndResponsesTestGroup() *TestGroup {
	tg := NewTestGroup("8.1.2.6", "Malformed Requests and Responses")

	tg.AddTestCase(NewTestCase(
		"Sends a HEADERS frame that contains the \"content-length\" header field which does not equal the sum of the DATA frame payload lengths",
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "POST"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair("content-length", "1"),
			}

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
		"the endpoint MUST respond with a stream error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "POST"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair("content-length", "1"),
			}

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
