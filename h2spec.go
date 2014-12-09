package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"net"
	"os"
	"strings"
	"time"
)

type TcpConn struct {
	conn   net.Conn
	dataCh chan []byte
	errCh  chan error
}

type Http2Conn struct {
	conn   net.Conn
	fr     *http2.Framer
	dataCh chan http2.Frame
	errCh  chan error
}

type Context struct {
	Port int
	Host string
}

func (ctx *Context) Authority() (authority string) {
	return fmt.Sprintf("%s:%d", ctx.Host, ctx.Port)
}

func Run(ctx *Context) {
	TestHttp2ConnectionPreface(ctx)
	TestFrameSize(ctx)
	TestHeaderCompressionAndDecompression(ctx)
	TestStreamStates(ctx)
	TestErrorHandling(ctx)
	TestData(ctx)
	TestHeaders(ctx)
	TestPriority(ctx)
	TestRstStream(ctx)
	TestSettings(ctx)
}

func CreateTcpConn(ctx *Context) *TcpConn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ctx.Host, ctx.Port))
	if err != nil {
		fmt.Println("Unable to connect to the target server.")
		os.Exit(1)
	}

	dataCh := make(chan []byte)
	errCh := make(chan error, 1)

	tcpConn := &TcpConn{
		conn:   conn,
		dataCh: dataCh,
		errCh:  errCh,
	}

	go func() {
		for {
			buf := make([]byte, 512)
			_, err := conn.Read(buf)
			dataCh <- buf
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	return tcpConn
}

func CreateHttp2Conn(ctx *Context, sn bool) *Http2Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ctx.Host, ctx.Port))
	if err != nil {
		fmt.Println("Unable to connect to the target server.")
		os.Exit(1)
	}

	fmt.Fprintf(conn, "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")

	if sn {
		done := false
		fr := http2.NewFramer(conn, conn)
		fr.WriteSettings()

		for {
			f, _ := fr.ReadFrame()
			switch f := f.(type) {
			case *http2.SettingsFrame:
				if f.IsAck() {
					done = true
				} else {
					fr.WriteSettingsAck()
				}
			default:
				done = true
			}

			if done {
				break
			}
		}
	}

	fr := http2.NewFramer(conn, conn)
	fr.AllowIllegalWrites = true
	dataCh := make(chan http2.Frame)
	errCh := make(chan error, 1)

	http2Conn := &Http2Conn{
		conn:   conn,
		fr:     fr,
		dataCh: dataCh,
		errCh:  errCh,
	}

	go func() {
		for {
			f, err := fr.ReadFrame()
			dataCh <- f
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	return http2Conn
}

func SetReadTimer(conn net.Conn, sec time.Duration) {
	now := time.Now()
	conn.SetReadDeadline(now.Add(time.Second * sec))
}

func PrintHeader(title string, i int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", i), title)
}

func PrintFooter() {
	fmt.Println("")
}

func PrintResult(result bool, desc string, msg string, i int) {
	var mark string
	indent := strings.Repeat("  ", i+1)
	if result {
		mark = "✓"
		fmt.Printf("%s\x1b[32m%s\x1b[0m \x1b[90m%s\x1b[0m\n", indent, mark, desc)
	} else {
		mark = "×"
		fmt.Printf("%s\x1b[31m%s %s\x1b[0m\n", indent, mark, desc)
		fmt.Printf("%s\x1b[31m  - %s\x1b[0m\n", indent, msg)
	}
}

func pair(name, value string) hpack.HeaderField {
	return hpack.HeaderField{Name: name, Value: value}
}
