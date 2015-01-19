package h2spec

import (
	"bytes"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"time"
)

func TestStreamStates(ctx *Context) {
	PrintHeader("5.1. Stream States", 0)
	TestStreamIdentifiers(ctx)
	PrintFooter()
}

func TestStreamIdentifiers(ctx *Context) {
	PrintHeader("5.1.1. Stream Identifiers", 1)
	msg := "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR."

	func(ctx *Context) {
		desc := "Sends even-numbered stream identifier"
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
		hp.StreamID = 2
		hp.EndStream = true
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
				}
			}
		}

		PrintResult(result, desc, msg, 1)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends stream identifier that is numerically smaller than previous"
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

		var hp1 http2.HeadersFrameParam
		hp1.StreamID = 5
		hp1.EndStream = true
		hp1.EndHeaders = true
		hp1.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp1)

		var hp2 http2.HeadersFrameParam
		hp2.StreamID = 3
		hp2.EndStream = true
		hp2.EndHeaders = true
		hp2.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp2)

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
				}
			}
		}

		PrintResult(result, desc, msg, 1)
	}(ctx)
}
