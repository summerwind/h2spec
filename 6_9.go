package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
)

func TestWindowUpdate(ctx *Context) {
	if !ctx.IsTarget("6.9") {
		return
	}

	PrintHeader("6.9. WINDOW_UPDATE", 0)

	func(ctx *Context) {
		desc := "Sends a WINDOW_UPDATE frame with an flow control window increment of 0"
		msg := "the endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		http2Conn.fr.WriteWindowUpdate(0, 0)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a WINDOW_UPDATE frame with an flow control window increment of 0 on a stream"
		msg := "the endpoint MUST respond with a stream error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)
		http2Conn.fr.WriteWindowUpdate(1, 0)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeProtocol {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a WINDOW_UPDATE frame with a length other than a multiple of 4 octets"
		msg := "the endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x03\x08\x00\x00\x00\x00\x00")
		fmt.Fprintf(http2Conn.conn, "\x00\x00\x01")

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				switch f.ErrCode {
				case http2.ErrCodeProtocol, http2.ErrCodeFrameSize:
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	TestFlowControlWindowLimits(ctx)
	TestInitialFlowControlWindowSize(ctx)

	PrintFooter()
}

func TestFlowControlWindowLimits(ctx *Context) {
	PrintHeader("6.9.1. The Flow Control Window", 1)

	func(ctx *Context) {
		desc := "Sends multiple WINDOW_UPDATE frames on a connection increasing the flow control window to above 2^31-1"
		msg := "the endpoint MUST respond with a stream error of type FLOW_CONTROL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		http2Conn.fr.WriteWindowUpdate(0, 2147483647)
		http2Conn.fr.WriteWindowUpdate(0, 2147483647)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFlowControl {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends multiple WINDOW_UPDATE frames on a stream increasing the flow control window to above 2^31-1"
		msg := "the endpoint MUST respond with a stream error of type FLOW_CONTROL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		hdrs := []hpack.HeaderField{
			pair(":method", "GET"),
			pair(":scheme", "http"),
			pair(":path", "/"),
			pair(":authority", ctx.Authority()),
		}

		var hp http2.HeadersFrameParam
		hp.StreamID = 1
		hp.EndStream = false
		hp.EndHeaders = true
		hp.BlockFragment = http2Conn.EncodeHeader(hdrs)
		http2Conn.fr.WriteHeaders(hp)

		http2Conn.fr.WriteWindowUpdate(1, 2147483647)
		http2Conn.fr.WriteWindowUpdate(1, 2147483647)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.RSTStreamFrame:
				if f.ErrCode == http2.ErrCodeFlowControl {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)
}

func TestInitialFlowControlWindowSize(ctx *Context) {
	PrintHeader("6.9.2. Initial Flow Control Window Size", 1)

	func(ctx *Context) {
		desc := "Sends a SETTINGS_INITIAL_WINDOW_SIZE settings with an exceeded maximum window size value"
		msg := "the endpoint MUST respond with a connection error of type FLOW_CONTROL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x04\x00\x00\x00\x00\x00")
		fmt.Fprintf(http2Conn.conn, "\x00\x04\x80\x00\x00\x00")

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFlowControl {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 1)
	}(ctx)
}
