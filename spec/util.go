package spec

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/summerwind/h2spec/config"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

var (
	gray   = color.New(color.FgHiBlack).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

// DummyString returns a dummy string with specified length.
func DummyString(len int) string {
	var buffer bytes.Buffer
	for i := 0; i < len; i++ {
		buffer.WriteString("x")
	}
	return buffer.String()
}

// DummyBytes returns a array of byte with specified length.
func DummyBytes(len int) []byte {
	var buffer bytes.Buffer
	for i := 0; i < len; i++ {
		buffer.WriteString("x")
	}
	return buffer.Bytes()
}

// HeaderField returns a header field of HPACK with specified
// name and value.
func HeaderField(name, value string) hpack.HeaderField {
	return hpack.HeaderField{Name: name, Value: value}
}

// CommonHeaders returns a array of header field of HPACK contained
// common http headers used in various test case.
func CommonHeaders(c *config.Config) []hpack.HeaderField {
	var scheme, authority string
	defaultPort := false

	if c.TLS {
		scheme = "https"
		if c.Port == 443 {
			defaultPort = true
		}
	} else {
		scheme = "http"
		if c.Port == 80 {
			defaultPort = true
		}
	}

	if defaultPort {
		authority = c.Host
	} else {
		authority = c.Addr()
	}

	return []hpack.HeaderField{
		HeaderField(":method", "GET"),
		HeaderField(":scheme", scheme),
		HeaderField(":path", c.Path),
		HeaderField(":authority", authority),
	}
}

// CommonRespHeaders returns a array of header field of HPACK contained
// common http headers used in various test case.
func CommonRespHeaders(c *config.Config) []hpack.HeaderField {
	return []hpack.HeaderField{
		HeaderField(":status", "200"),
		HeaderField("access-control-allow-origin", "*"),
	}
}

// DummyHeaders returns a array of header field of HPACK contained
// dummy string values.
func DummyHeaders(c *config.Config, len int) []hpack.HeaderField {
	headers := make([]hpack.HeaderField, 0, len)
	dummy := DummyString(c.MaxHeaderLen)

	for i := 0; i < len; i++ {
		name := fmt.Sprintf("x-dummy%d", i)
		headers = append(headers, HeaderField(name, dummy))
	}

	return headers
}

func DummyRespHeaders(c *config.Config, len int) []hpack.HeaderField {
	headers := make([]hpack.HeaderField, 0, len)
	dummy := DummyString(c.MaxHeaderLen)

	for i := 0; i < len; i++ {
		name := fmt.Sprintf("x-dummy%d", i)
		headers = append(headers, HeaderField(name, dummy))
	}

	return headers
}

// ServerDataLength returns the total length of the DATA frame of /.
func ServerDataLength(c *config.Config) (int, error) {
	conn, err := Dial(c)
	if err != nil {
		return 0, err
	}

	err = conn.Handshake()
	if err != nil {
		return 0, err
	}

	headers := CommonHeaders(c)
	hp := http2.HeadersFrameParam{
		StreamID:      1,
		EndStream:     true,
		EndHeaders:    true,
		BlockFragment: conn.EncodeHeaders(headers),
	}
	conn.WriteHeaders(hp)

	len := 0
	done := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case DataFrameEvent:
			len += int(event.Header().Length)
			done = event.StreamEnded()
		case HeadersFrameEvent:
			done = event.StreamEnded()
		}

		if done {
			break
		}
	}

	if !done {
		return 0, errors.New("Unable to get server data length")
	}

	return len, nil
}
