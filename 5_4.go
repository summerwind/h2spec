package h2spec

import (
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"

	"golang.org/x/net/http2"
)

func ErrorHandlingTestGroup(ctx *Context) *TestGroup {
	if !ctx.Strict {
		return nil
	}

	tg := NewTestGroup("5.4", "Error Handling")
	tg.AddTestGroup(ConnectionErrorHandlingTestGroup(ctx))

	return tg
}

func ConnectionErrorHandlingTestGroup(ctx *Context) *TestGroup {
	tg := NewTestGroup("5.4.1", "Connection Error Handling")

	tg.AddTestCase(NewTestCase(
		"Raise a connection error",
		"After sending the GOAWAY frame, the endpoint MUST close the TCP connection.",
		func(ctx *Context) (pass bool, expected []Result, actual Result) {
			pass = false
			expected = []Result{
				&ResultConnectionClose{},
			}
			goaway := false
			closed := false

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			// PING frame with invalid stream ID
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x08\x06\x00\x00\x00\x00\x03")
			fmt.Fprintf(http2Conn.conn, "\x00\x00\x00\x00\x00\x00\x00\x00")

		loop:
			for {
				f, err := http2Conn.ReadFrame(ctx.Timeout)
				if err != nil {
					opErr, ok := err.(*net.OpError)
					if err == io.EOF || (ok && opErr.Err == syscall.ECONNRESET) {
						closed = true
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
				case *http2.GoAwayFrame:
					if f.ErrCode == http2.ErrCodeProtocol {
						goaway = true
					}
				default:
					actual = CreateResultFrame(f)
				}
			}

			if goaway && closed && actual == nil {
				actual = &ResultConnectionClose{}
				pass = true
			} else {
				actual = &ResultError{
					errors.New("Connection closed, but did not receive a GOAWAY Frame."),
				}
			}

			return pass, expected, actual
		},
	))

	return tg
}
