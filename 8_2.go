package h2spec

import (
	"bytes"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"time"
)

func TestServerPush(ctx *Context) {
	if !ctx.IsTarget("8.2") {
		return
	}

	PrintHeader("8.2. Server Push", 0)

	func(ctx *Context) {
		desc := "Sends a PUSH_PROMISE frame"
		msg := "the endpoint MUST treat the receipt of a PUSH_PROMISE frame as a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		var buf bytes.Buffer
		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
		}
		enc := hpack.NewEncoder(&buf)
		for _, hf := range hdrs {
			_ = enc.WriteField(hf)
		}

		var pp http2.PushPromiseParam
		pp.StreamID = 1
		pp.PromiseID = 3
		pp.EndHeaders = true
		pp.BlockFragment = buf.Bytes()
		http2Conn.fr.WritePushPromise(pp)

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				switch f.ErrCode {
				case http2.ErrCodeProtocol:
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
