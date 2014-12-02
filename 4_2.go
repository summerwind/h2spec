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
		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf := f.(*http2.GoAwayFrame)
				if gf != nil {
					if gf.ErrCode == http2.ErrCodeFrameSize {
						result = true
						break loop
					}
				}
				break
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
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
		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf := f.(*http2.GoAwayFrame)
				if gf != nil {
					if gf.ErrCode == http2.ErrCodeFrameSize {
						result = true
						break loop
					}
				}
				break
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
