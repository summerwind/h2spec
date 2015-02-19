package h2spec

import (
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"io"
	"net"
	"syscall"
)

func ExtendingHttp2TestGroup() *TestGroup {
	tg := NewTestGroup("5.5", "Extending HTTP/2")

	tg.AddTestCase(NewTestCase(
		"Sends an unknown extension frame",
		"The endpoint MUST discard frames that have unknown or unsupported types",
		func(ctx *Context) (expected []Result, actual Result) {
			expected = []Result{
				&ResultFrame{http2.FramePing, http2.FlagPingAck, ErrCodeDefault},
			}

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
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						actual = &ResultConnectionClose{}
					} else if err == TIMEOUT {
						if actual == nil {
							actual = &ResultTestTimeout{}
						}
					} else {
						actual = &ResultError{err}
					}
					break loop
				}

				switch f := f.(type) {
				case *http2.PingFrame:
					actual = &ResultFrame{f.Header().Type, f.Header().Flags, ErrCodeDefault}
					if f.FrameHeader.Flags.Has(http2.FlagPingAck) {
						break loop
					}
				default:
					actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
				}
			}

			return expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends an unknown extension frame in the middle of a header block",
		"The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := []hpack.HeaderField{
				pair(":method", "GET"),
				pair(":scheme", "http"),
				pair(":path", "/"),
				pair(":authority", ctx.Authority()),
				pair("x-dummy1", dummyData(10000)),
				pair("x-dummy2", dummyData(10000)),
			}

			blockFragment := http2Conn.EncodeHeader(hdrs)

			var hp http2.HeadersFrameParam
			hp.StreamID = 1
			hp.EndStream = true
			hp.EndHeaders = false
			hp.BlockFragment = blockFragment[0:16384]
			http2Conn.fr.WriteHeaders(hp)

			http2Conn.fr.WriteRawFrame(0xFF, 0x01, 0, []byte("unknown"))

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
