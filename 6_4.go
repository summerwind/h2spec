package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func TestRstStream(ctx *Context) {
	if !ctx.IsTarget("6.4") {
		return
	}

	PrintHeader("6.4. RST_STREAM", 0)

	func(ctx *Context) {
		desc := "Sends a RST_STREAM frame with 0x0 stream identifier"
		msg := "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		http2Conn.fr.WriteRSTStream(0, http2.ErrCodeCancel)

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
		desc := "Sends a RST_STREAM frame on a idle stream"
		msg := "The endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		http2Conn.fr.WriteRSTStream(1, http2.ErrCodeCancel)

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
		desc := "Sends a RST_STREAM frame with a length other than 4 octets"
		msg := "The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

		// RST_STREAM Frame
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x03\x03\x00\x00\x00\x00\x01")
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x00")

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
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
