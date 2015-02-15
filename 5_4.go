package h2spec

import (
	"errors"
	"fmt"
	"github.com/bradfitz/http2"
	"io"
)

func ErrorHandlingTestGroup() *TestGroup {
	tg := NewTestGroup("5.4", "Error Handling")

	tg.AddTestGroup(ConnectionErrorHandlingTestGroup())

	return tg
}

func ConnectionErrorHandlingTestGroup() *TestGroup {
	tg := NewTestGroup("5.4.1", "Connection Error Handling")

	tg.AddTestCase(NewTestCase(
		"Receives a GOAWAY frame",
		"After sending the GOAWAY frame, the endpoint MUST close the TCP connection.",
		func(ctx *Context) (expected []Result, actual Result) {
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
					if err == io.EOF {
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
					actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
				}
			}

			if goaway && closed && actual == nil {
				actual = &ResultConnectionClose{}
			} else {
				actual = &ResultError{
					errors.New("Connection closed, but did not receive a GOAWAY Frame."),
				}
			}

			return expected, actual
		},
	))

	return tg
}
