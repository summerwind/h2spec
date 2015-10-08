package h2spec

import (
	"golang.org/x/net/http2"
	"io"
	"net"
	"syscall"
)

func ExtendingHttp2TestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("5.5", "Extending HTTP/2")

	tg.AddTestCase(NewTestCase(
		"Sends an unknown extension frame",
		"The endpoint MUST discard frames that have unknown or unsupported types",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultFrame{LengthDefault, http2.FramePing, http2.FlagPingAck, ErrCodeDefault},
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
						rf, ok := actual.(*ResultFrame)
						if actual == nil || (ok && rf.Type != http2.FrameGoAway) {
							actual = &ResultConnectionClose{}
						}
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
					actual = CreateResultFrame(f)
					if f.FrameHeader.Flags.Has(http2.FlagPingAck) {
						pass = true
						break loop
					}
				default:
					actual = CreateResultFrame(f)
				}
			}

			return pass, expected, actual
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends an unknown extension frame in the middle of a header block",
		"The endpoint MUST treat as a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			hdrs := commonHeaderFields(ctx)
			hdrs = append(hdrs, pair("x-dummy1", dummyData(10000)))
			hdrs = append(hdrs, pair("x-dummy2", dummyData(10000)))
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
