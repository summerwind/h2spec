package h2spec

import (
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"net"
	"syscall"
)

func PingTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("6.7", "PING")

	tg.AddTestCase(NewTestCase(
		"Sends a PING frame",
		"The endpoint MUST sends a PING frame with ACK, with an identical payload.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultFrame{8, http2.FramePing, http2.FlagPingAck, ErrCodeDefault},
			}

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

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
					if f.FrameHeader.Flags.Has(http2.FlagPingAck) && f.Data == data {
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
		"Sends a PING frame with the stream identifier that is not 0x0",
		"The endpoint MUST respond with a connection error of type PROTOCOL_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x08\x06\x00\x00\x00\x00\x03")
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x00\x00\x00\x00\x00")

			actualCodes := []http2.ErrCode{http2.ErrCodeProtocol}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	tg.AddTestCase(NewTestCase(
		"Sends a PING frame with a length field value other than 8",
		"The endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x06\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x00\x00\x00")

			actualCodes := []http2.ErrCode{http2.ErrCodeFrameSize}
			return TestConnectionError(ctx, http2Conn, actualCodes)
		},
	))

	return tg
}
