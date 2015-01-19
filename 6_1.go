package h2spec

import (
	"bytes"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"time"
)

func TestData(ctx *Context) {
	PrintHeader("6.1. DATA", 0)

	func(ctx *Context) {
		desc := "Sends a DATA frame with 0x0 stream identifier"
		msg := "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		http2Conn.fr.WriteData(0, true, []byte("test"))
		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf, ok := f.(*http2.GoAwayFrame)
				if ok {
					if gf.ErrCode == http2.ErrCodeProtocol {
						result = true
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
		desc := "Sends a DATA frame on the stream that is not opend"
		msg := "The endpoint MUST respond with a stream error of type STREAM_CLOSED."
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

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf, ok := f.(*http2.RSTStreamFrame)
				if ok {
					if gf.ErrCode == http2.ErrCodeStreamClosed {
						result = true
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

	// bradfitz/http2 does not support padding of DATA frame.
	/*
		func(ctx *Context) {
			desc := "Sends a DATA frame with invalid pad length"
			msg := "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR."
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

			timeCh := time.After(3 * time.Second)

		loop:
			for {
				select {
				case f := <-http2Conn.dataCh:
					gf, ok := f.(*http2.GoAwayFrame)
					if ok {
						if gf.ErrCode == http2.ErrCodeProtocol {
							result = true
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
	*/

	PrintFooter()
}
