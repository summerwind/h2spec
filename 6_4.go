package h2spec

import (
	"github.com/bradfitz/http2"
	"time"
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

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
