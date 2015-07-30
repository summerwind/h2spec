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
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			// Literal Header Field with Incremental Indexing without Length and String segment
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x01\x01\x05\x00\x00\x00\x01\x40")

			actualCodes := []http2.ErrCode{http2.ErrCodeCompression}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends Dynamic Table Size Update (RFC 7541, 6.3)",
		"The endpoint must accept Dynamic Table Size Update",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			// 2 Dynamic Table Size Updates, 0 and 4096.
			blockFragment := []byte{0x20, 0x3f, 0xe1, 0x1f}
			blockFragment = append(blockFragment, http2Conn.EncodeHeader(hdrs)...)
			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = blockFragment
			http2Conn.fr.WriteHeaders(hp)

			return TestStreamClose(ctx, http2Conn)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Encodes Dynamic Table Size Update (RFC 7541, 6.3) after common header fields",
		"The endpoint MUST terminate the connection with a connection error of type COMPRESSION_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)

			blockFragment := http2Conn.EncodeHeader(hdrs)
			// append 2 Dynamic Table Size Updates, 0 and
			// 4096.  this is illegal, since RFC 7541,
			// section 4.2 says that dynamic table size
			// update MUST occur at the beginning of the
			// first header block following the changes to
			// the dynamic table size.
			blockFragment = append(blockFragment, 0x20, 0x3f, 0xe1, 0x1f)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = true
			hp.BlockFragment = blockFragment
			http2Conn.fr.WriteHeaders(hp)

			actualCodes := []http2.ErrCode{http2.ErrCodeCompression}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
