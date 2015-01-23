package h2spec

import (
	"fmt"
	"io"
	"time"
)

func TestHttp2ConnectionPreface(ctx *Context) {
	if !ctx.IsTarget("3.5") {
		return
	}

	PrintHeader("3.5. HTTP/2 Connection Preface", 0)

	func(ctx *Context) {
		desc := "Sends invalid connection preface"
		msg := "The endpoint MUST terminate the TCP connection."
		result := false

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
					result = true
					break loop
				}
			case <-timeCh:
				break loop
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	PrintFooter()
}
