package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
)

func HeaderCompressionAndDecompressionTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("4.3", "Header Compression and Decompression")

	tg.AddTestCase(NewTestCase(
		"Sends invalid header block fragment",
		"The endpoint MUST terminate the connection with a connection error of type COMPRESSION_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			// Literal Header Field with Incremental Indexing without Length and String segment
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x01\x01\x05\x00\x00\x00\x01\x40")

			actualCodes := []http2.ErrCode{http2.ErrCodeCompression}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
