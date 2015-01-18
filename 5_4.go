package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"io"
	"time"
)

func TestErrorHandling(ctx *Context) {
	PrintHeader("5.4. Error Handling", 0)
	TestConnectionErrorHandling(ctx)
	PrintFooter()
}

func TestConnectionErrorHandling(ctx *Context) {
	PrintHeader("5.4.1. Connection Error Handling", 1)

	func(ctx *Context) {
		desc := "Receives a GOAWAY frame"
		msg := "After sending the GOAWAY frame, the endpoint MUST close the TCP connection."
		gfResult := false
		closeResult := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		// PING frame with invalid stream ID
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x08\x06\x00\x00\x00\x00\x03")
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x00\x00\x00\x00\x00")

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				if err == io.EOF {
					closeResult = true
				}
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					gfResult = true
				}
			}
		}

		PrintResult(gfResult && closeResult, desc, msg, 1)
	}(ctx)
}
