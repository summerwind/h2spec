package h2spec

import (
	"fmt"
	"io"
	"time"
)

func Http2ConnectionPrefaceTestGroup() *TestGroup {
	tg := NewTestGroup("3.5", "HTTP/2 Connection Preface")

	tg.AddTestCase(NewTestCase(
		"Sends invalid connection preface",
		"The endpoint MUST terminate the TCP connection.",
		func(ctx *Context) (expected []Result, actual Result) {
			expected = []Result{
				&ResultConnectionClose{},
			}

			tcpConn := CreateTcpConn(ctx)
			defer tcpConn.conn.Close()

			fmt.Fprintf(tcpConn.conn, "INVALID CONNECTION PREFACE")
			timeCh := time.After(ctx.Timeout)

		loop:
			for {
				select {
				case <-tcpConn.dataCh:
					break
				case err := <-tcpConn.errCh:
					if err == io.EOF {
						actual = &ResultConnectionClose{}
					} else {
						actual = &ResultError{err}
					}
					break loop
				case <-timeCh:
					actual = &ResultTestTimeout{}
					break loop
				}
			}

			return expected, actual
		},
	))

	return tg
}
