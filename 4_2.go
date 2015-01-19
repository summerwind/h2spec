package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"time"
)

func TestFrameSize(ctx *Context) {
	PrintHeader("4.2. Frame Size", 0)
	msg := "The endpoint MUST send a FRAME_SIZE_ERROR error."

	func(ctx *Context) {
		desc := "Sends too small size frame"
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x08\x00\x00\x00\x00\x00")

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFrameSize {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends too large size frame"
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x0f\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFrameSize {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
