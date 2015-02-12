package h2spec

import (
	"github.com/bradfitz/http2"
	"io"
)

func TestExtendingHTTP2(ctx *Context) {
	if !ctx.IsTarget("5.5") {
		return
	}

	PrintHeader("5.5. Extending HTTP/2", 0)
	TestUnknownFrames(ctx)
	PrintFooter()
}

func TestUnknownFrames(ctx *Context) {

	func(ctx *Context) {
		desc := "Sends an unknown frame type (0xFF)"
		msg := "the endpoint must ignore unknown frame types."
		result := false
		closeResult := true

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		// Write a frame of type 0xFF, which isn't yet defined
		// as an extension frame. This should be ignored; no GOAWAY,
		// RST_STREAM or closing the connection should occur
		http2Conn.fr.WriteRawFrame(0xFF, 0x00, 0, []byte("unknown"))

		// Now send a normal PING frame, and if this is processed
		// without error, then the preceeding unknown frame must have
		// been processed and ignored.
		data := [8]byte{'h', '2', 's', 'p', 'e', 'c'}
		http2Conn.fr.WritePing(false, data)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				if err == io.EOF {
					closeResult = false
				}
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = false
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = false
				}
			case *http2.PingFrame:
				if f.FrameHeader.Flags.Has(http2.FlagPingAck) {
					result = true
					break loop
				}
			}
		}

		PrintResult(result && closeResult, desc, msg, 1)
	}(ctx)
}
