package h2spec

import (
	"bytes"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"time"
)

func GetDummyData(num int) string {
	var data string
	for i := 0; i < num; i++ {
		data += "x"
	}
	return data
}

func TestContinuation(ctx *Context) {
	PrintHeader("6.10. CONTINUATION", 0)

	func(ctx *Context) {
		desc := "Sends a CONTINUATION frame"
		msg := "The endpoint must accept the frame."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("x-dummy1", GetDummyData(10000)),
			pair("x-dummy2", GetDummyData(10000)),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var blockFragment = buf.Bytes()

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = false
		hp.BlockFragment = blockFragment[0:16384]
		http2Conn.fr.WriteHeaders(hp)

		http2Conn.fr.WriteContinuation(1, true, blockFragment[16384:])

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				_, ok := f.(*http2.HeadersFrame)
				if ok {
					result = true
					break loop
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends multiple CONTINUATION frames"
		msg := "The endpoint must accept the frames."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("x-dummy1", GetDummyData(10000)),
			pair("x-dummy2", GetDummyData(10000)),
			pair("x-dummy3", GetDummyData(10000)),
			pair("x-dummy4", GetDummyData(10000)),
			pair("x-dummy5", GetDummyData(10000)),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var blockFragment = buf.Bytes()

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = false
		hp.BlockFragment = blockFragment[0:16384]
		http2Conn.fr.WriteHeaders(hp)

		http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])
		http2Conn.fr.WriteContinuation(1, true, blockFragment[32767:])

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				_, ok := f.(*http2.HeadersFrame)
				if ok {
					result = true
					break loop
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a CONTINUATION frame followed by any frame other than CONTINUATION"
		msg := "The endpoint MUST treat as a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("x-dummy1", GetDummyData(10000)),
			pair("x-dummy2", GetDummyData(10000)),
			pair("x-dummy3", GetDummyData(10000)),
			pair("x-dummy4", GetDummyData(10000)),
			pair("x-dummy5", GetDummyData(10000)),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var blockFragment = buf.Bytes()

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = false
		hp.BlockFragment = blockFragment[0:16384]
		http2Conn.fr.WriteHeaders(hp)

		http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])

		http2Conn.fr.WriteData(1, true, []byte("test"))

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf, ok := f.(*http2.GoAwayFrame)
				if ok {
					if gf.ErrCode == http2.ErrCodeProtocol {
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

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a CONTINUATION frame followed by a frame on a different stream"
		msg := "The endpoint MUST treat as a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("x-dummy1", GetDummyData(10000)),
			pair("x-dummy2", GetDummyData(10000)),
			pair("x-dummy3", GetDummyData(10000)),
			pair("x-dummy4", GetDummyData(10000)),
			pair("x-dummy5", GetDummyData(10000)),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var blockFragment = buf.Bytes()

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = false
		hp.BlockFragment = blockFragment[0:16384]
		http2Conn.fr.WriteHeaders(hp)

		http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])
		http2Conn.fr.WriteContinuation(3, true, blockFragment[32767:])

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf, ok := f.(*http2.GoAwayFrame)
				if ok {
					if gf.ErrCode == http2.ErrCodeProtocol {
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

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a CONTINUATION frame with the stream identifier that is 0x0"
		msg := "The endpoint MUST treat as a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
			pair("x-dummy1", GetDummyData(10000)),
			pair("x-dummy2", GetDummyData(10000)),
			pair("x-dummy3", GetDummyData(10000)),
			pair("x-dummy4", GetDummyData(10000)),
			pair("x-dummy5", GetDummyData(10000)),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var blockFragment = buf.Bytes()

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = true
		hp.EndHeaders = false
		hp.BlockFragment = blockFragment[0:16384]
		http2Conn.fr.WriteHeaders(hp)

		http2Conn.fr.WriteContinuation(1, false, blockFragment[16384:32767])
		http2Conn.fr.WriteContinuation(0, true, blockFragment[32767:])

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf, ok := f.(*http2.GoAwayFrame)
				if ok {
					if gf.ErrCode == http2.ErrCodeProtocol {
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

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a CONTINUATION frame after the frame other than HEADERS, PUSH_PROMISE or CONTINUATION"
		msg := "The endpoint MUST treat as a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
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

		http2Conn.fr.WriteData(1, true, []byte("test"))

		http2Conn.fr.WriteContinuation(1, true, buf.Bytes())

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf, ok := f.(*http2.GoAwayFrame)
				if ok {
					if gf.ErrCode == http2.ErrCodeProtocol {
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

		PrintResult(result, desc, msg, 0)
	}(ctx)
}
