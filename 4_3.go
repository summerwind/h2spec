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
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x14\x01\x05\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f:= f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeCompression {
					result = true
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
