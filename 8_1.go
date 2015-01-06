package h2spec

import (
	"bytes"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"time"
)

func TestHTTPRequestResponseExchange(ctx *Context) {
	PrintHeader("8.1. HTTP Request/Response Exchange", 0)

	TestHTTPHeaderFields(ctx)

	PrintFooter()
}

func TestHTTPHeaderFields(ctx *Context) {
	PrintHeader("8.1.2. HTTP Header Fields", 1)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the header field name in uppercase letters"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("X-TEST", "test"),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 1)
	}(ctx)

	TestPseudoHeaderFields(ctx)
	TestConnectionSpecificHeaderFields(ctx)
	TestRequestPseudoHeaderFields(ctx)
	TestMalformedRequestsAndResponses(ctx)
}

func TestPseudoHeaderFields(ctx *Context) {
	PrintHeader("8.1.2.1. Pseudo-Header Fields", 2)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the pseudo-header field defined for response"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair(":status", "200"),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the invalid pseudo-header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair(":test", "test"),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains a pseudo-header field that appears in a header block after a regular header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair("x-test", "test"),
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}

func TestConnectionSpecificHeaderFields(ctx *Context) {
	PrintHeader("8.1.2.2. Connection-Specific Header Fields", 2)

	// 8.1.2.3. CONNECT リクエスト (8.3節) である場合を除き、":method"、":scheme"、そして ":path" 擬似ヘッダーフィールドに有効な値を1つ含まなければなりません (MUST)。これらの擬似ヘッダーフィールドが省略された HTTP リクエストは不正な形式 (8.1.2.6節) です。
	// 8.1.2.6. ボディを構成する DATA フレームペイロードの長さの合計が "content-length" ヘッダーフィールドの値と等しくない場合、リクエストやレスポンスは不正な形式になります。

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the connection-specific header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("connection", "keep-alive"),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains the TE header field that contain any value other than \"trailers\""
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("trailers", "test"),
			pair("te", "trailers, deflate"),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}

func TestRequestPseudoHeaderFields(ctx *Context) {
	PrintHeader("8.1.2.3. Request Pseudo-Header Fields", 2)

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that is omitted mandatory pseudo-header fields"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":authority", ctx.Authority()),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}

func TestMalformedRequestsAndResponses(ctx *Context) {
	PrintHeader("8.1.2.6. Malformed Requests and Responses", 2)

	// 8.1.2.6. ボディを構成する DATA フレームペイロードの長さの合計が "content-length" ヘッダーフィールドの値と等しくない場合、リクエストやレスポンスは不正な形式になります。

	func(ctx *Context) {
		desc := "Sends a HEADERS frame that contains invalid \"content-length\" header field"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "POST"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("content-length", "1"),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)
		http2Conn.fr.WriteData(1, true, []byte("test"))

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				rf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if rf.ErrCode == http2.ErrCodeProtocol {
						result = true
						break loop
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 2)
	}(ctx)
}
