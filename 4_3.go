package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"time"
)

func TestHeaderCompressionAndDecompression(ctx *Context) {
	PrintHeader("4.3. Header Compression and Decompression", 0)

	func(ctx *Context) {
		desc := "Sends invalid header block fragment"
		msg := "The endpoint MUST terminate the connection with a connection error of type COMPRESSION_ERROR."
		rfResult := false
		gfResult := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x14\x01\x05\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
		timeCh := time.After(3 * time.Second)

	loop:
		for {
			select {
			case f := <-http2Conn.dataCh:
				switch frame := f.(type) {
				case *http2.RSTStreamFrame:
					if frame.ErrCode == http2.ErrCodeCompression {
						rfResult = true
					}
				case *http2.GoAwayFrame:
					if frame.ErrCode == http2.ErrCodeCompression {
						gfResult = true
					}
				}
			case <-http2Conn.errCh:
				break loop
			case <-timeCh:
				break loop
			}
		}

		PrintResult(rfResult && gfResult, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
