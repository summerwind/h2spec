package spec

import (
	"bytes"
	"fmt"

	"github.com/fatih/color"
	"github.com/summerwind/h2spec/config"

	"golang.org/x/net/http2/hpack"
)

var (
	gray   = color.New(color.FgHiBlack).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

func DummyString(len int) string {
	var buffer bytes.Buffer
	for i := 0; i < len; i++ {
		buffer.WriteString("x")
	}
	return buffer.String()
}

func DummyBytes(len int) []byte {
	var buffer bytes.Buffer
	for i := 0; i < len; i++ {
		buffer.WriteString("x")
	}
	return buffer.Bytes()
}

func HeaderField(name, value string) hpack.HeaderField {
	return hpack.HeaderField{Name: name, Value: value}
}

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
		HeaderField(":path", "/"),
		HeaderField(":authority", authority),
	}
}

func DummyHeaders(c *config.Config, len int) []hpack.HeaderField {
	headers := make([]hpack.HeaderField, len)
	dummy := DummyString(c.MaxHeaderLen)

	for i := 0; i < len; i++ {
		name := fmt.Sprintf("x-dummy%d", i)
		headers = append(headers, HeaderField(name, dummy))
	}

	return headers
}
