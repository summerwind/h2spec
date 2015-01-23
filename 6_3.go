package h2spec

import (
	"bytes"
	"fmt"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func TestPriority(ctx *Context) {
	if !ctx.IsTarget("6.3") {
		return
	}

	PrintHeader("6.3. PRIORITY", 0)

	func(ctx *Context) {
		desc := "Sends a PRIORITY frame with 0x0 stream identifier"
		msg := "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
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
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		// PRIORITY Frame
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x05\x02\x00\x00\x00\x00\x00")
		fmt.Fprintf(http2Conn.conn, "\x80\x00\x00\x01\x0a")

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
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

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a PRIORITY frame with a length other than 5 octets"
		msg := "The endpoint MUST respond with a stream error of type FRAME_SIZE_ERROR."
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
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = buf.Bytes()
		http2Conn.fr.WriteHeaders(hp)

		// PRIORITY Frame
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x04\x02\x00\x00\x00\x00\x01")
		fmt.Fprintf(http2Conn.conn, "\x80\x00\x00\x01")

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFrameSize {
					result = true
				}
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeFrameSize {
					result = true
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
