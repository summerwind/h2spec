package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"time"
)

func TestPing(ctx *Context) {
	if !ctx.IsTarget("6.7") {
		return
	}

	PrintHeader("6.7. PING", 0)

	func(ctx *Context) {
		desc := "Sends a PING frame"
		msg := "the endpoint MUST sends a PING frame with ACK."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		data := [8]byte{'h', '2', 's', 'p', 'e', 'c'}
		http2Conn.fr.WritePing(false, data)

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.PingFrame:
				if f.FrameHeader.Flags.Has(http2.FlagPingAck) {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a PING frame with the stream identifier that is not 0x0"
		msg := "the endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x08\x06\x00\x00\x00\x00\x03")
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

	func(ctx *Context) {
		desc := "Sends a PING frame with a length field value other than 8"
		msg := "the endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x06\x00\x00\x00\x00\x00")
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x00\x00\x00")

	loop:
		for {
			f, err := http2Conn.ReadFrame(3 * time.Second)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFrameSize {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
