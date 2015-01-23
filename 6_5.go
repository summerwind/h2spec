package h2spec

import (
	"fmt"
	"github.com/bradfitz/http2"
)

func TestSettings(ctx *Context) {
	if !ctx.IsTarget("6.5") {
		return
	}

	PrintHeader("6.5. SETTINGS", 0)

	func(ctx *Context) {
		desc := "Sends a SETTINGS frame"
		msg := "the endpoint MUST sends a SETTINGS frame with ACK."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		settings := http2.Setting{http2.SettingMaxConcurrentStreams, 100}
		http2Conn.fr.WriteSettings(settings)

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.SettingsFrame:
				if f.IsAck() {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a SETTINGS frame that is not a zero-length with ACK flag"
		msg := "the endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x01\x04\x01\x00\x00\x00\x00\x00")

	loop:
		for {
			f, err := http2Conn.ReadFrame(ctx.Timeout)
			if err != nil {
				break loop
			}
			switch f := f.(type) {
			case *http2.GoAwayFrame:
				if f.ErrCode == http2.ErrCodeFrameSize {
					result = true
					break loop
				}
			}
		}

		PrintResult(result, desc, msg, 0)
	}(ctx)

	func(ctx *Context) {
		desc := "Sends a SETTINGS frame with the stream identifier that is not 0x0"
		msg := "the endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x04\x00\x00\x00\x00\x03")
		fmt.Fprintf(http2Conn.conn, "\x00\x03\x00\x00\x00\x64")

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
		desc := "Sends a SETTINGS frame with a length other than a multiple of 6 octets"
		msg := "the endpoint MUST respond with a connection error of type FRAME_SIZE_ERROR."
		result := false

		http2Conn := CreateHttp2Conn(ctx, true)
		defer http2Conn.conn.Close()

		fmt.Fprintf(http2Conn.conn, "\x00\x00\x02\x04\x00\x00\x00\x00\x00")
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

	TestDefinedSettingsParameters(ctx)
	PrintFooter()
}

func TestDefinedSettingsParameters(ctx *Context) {
	PrintHeader("6.5.2. Defined SETTINGS Parameters", 1)

	func(ctx *Context) {
		PrintHeader("SETTINGS_ENABLE_PUSH (0x2)", 2)

		func(ctx *Context) {
			desc := "Sends the value other than 0 or 1"
			msg := "the endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
			result := false

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x04\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x00\x02\x00\x00\x00\x02")

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

			PrintResult(result, desc, msg, 2)
		}(ctx)
	}(ctx)

	func(ctx *Context) {
		PrintHeader("SETTINGS_INITIAL_WINDOW_SIZE (0x4)", 2)

		func(ctx *Context) {
			desc := "Sends the value above the maximum flow control window size"
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

			PrintResult(result, desc, msg, 2)
		}(ctx)
	}(ctx)

	func(ctx *Context) {
		PrintHeader("SETTINGS_MAX_FRAME_SIZE (0x5)", 2)

		func(ctx *Context) {
			desc := "Sends the value below the initial value"
			msg := "the endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
			result := false

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x04\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x00\x05\x00\x00\x3f\xff")

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

			PrintResult(result, desc, msg, 2)
		}(ctx)

		func(ctx *Context) {
			desc := "Sends the value above the maximum allowed frame size"
			msg := "the endpoint MUST respond with a connection error of type PROTOCOL_ERROR."
			result := false

			http2Conn := CreateHttp2Conn(ctx, true)
			defer http2Conn.conn.Close()

			fmt.Fprintf(http2Conn.conn, "\x00\x00\x06\x04\x00\x00\x00\x00\x00")
			fmt.Fprintf(http2Conn.conn, "\x00\x05\x01\x00\x00\x00")

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

			PrintResult(result, desc, msg, 2)
		}(ctx)
	}(ctx)
}
