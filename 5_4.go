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

		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				gf, ok := f.(*http2.GoAwayFrame)
				if ok {
					if gf.ErrCode == http2.ErrCodeProtocol {
						gfResult = true
					}
				}
			case err := <-http2Conn.errCh:
				if err == io.EOF {
					closeResult = true
				}
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(gfResult && closeResult, desc, msg, 1)
	}(ctx)
}
