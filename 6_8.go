package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"time"
)

func TestGoaway(ctx *Context) {
	if !ctx.IsTarget("6.8") {
		return
	}

	PrintHeader("6.8. GOAWAY", 0)

	func(ctx *Context) {
		desc := "Sends a GOAWAY frame with the stream identifier that is not 0x0"
		msg := "the endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x08\x07\x00\x00\x00\x00\x03")
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x00\x00\x00\x00\x00")

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
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
