package spec

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
)

const DefaultWindowSize = 65535

type Conn struct {
	net.Conn

	Settings map[http2.SettingID]uint32
	Timeout  time.Duration
	Verbose  bool
	Closed   bool

	WindowUpdate bool
	WindowSize   map[uint32]int

	framer     *http2.Framer
	encoder    *hpack.Encoder
	encoderBuf *bytes.Buffer

	debugFramer    *http2.Framer
	debugFramerBuf *bytes.Buffer
}

func Dial(c *config.Config) (*Conn, error) {
	var baseConn net.Conn
	var err error

	if c.TLS {
		dialer := &net.Dialer{}
		dialer.Timeout = c.Timeout

		tlsConn, err := tls.DialWithDialer(dialer, "tcp", c.Addr(), c.TLSConfig())
		if err != nil {
			return nil, err
		}

		cs := tlsConn.ConnectionState()
		if !cs.NegotiatedProtocolIsMutual {
			return nil, errors.New("Protocol negotiation failed")
		}

		baseConn = tlsConn
	} else {
		baseConn, err = net.DialTimeout("tcp", c.Addr(), c.Timeout)
		if err != nil {
			return nil, err
		}
	}

	settings := map[http2.SettingID]uint32{}

	framer := http2.NewFramer(baseConn, baseConn)
	framer.AllowIllegalWrites = true
	framer.AllowIllegalReads = true

	var encoderBuf bytes.Buffer
	encoder := hpack.NewEncoder(&encoderBuf)

	conn := Conn{
		Conn:     baseConn,
		Settings: settings,
		Timeout:  c.Timeout,
		Verbose:  c.Verbose,
		Closed:   false,

		WindowUpdate: true,
		WindowSize:   map[uint32]int{0: DefaultWindowSize},

		framer:     framer,
		encoder:    encoder,
		encoderBuf: &encoderBuf,
	}

	if conn.Verbose {
		conn.debugFramerBuf = new(bytes.Buffer)
		conn.debugFramer = http2.NewFramer(conn.debugFramerBuf, conn.debugFramerBuf)
		conn.debugFramer.AllowIllegalWrites = true
		conn.debugFramer.AllowIllegalReads = true
	}

	return &conn, nil
}

func (conn *Conn) Handshake() error {
	done := make(chan error)

	fmt.Fprintf(conn, "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")

	go func() {
		local := false
		remote := false

		setting := http2.Setting{
			ID:  http2.SettingInitialWindowSize,
			Val: DefaultWindowSize,
		}
		conn.WriteSettings(setting)

		for !(local && remote) {
			f, err := conn.framer.ReadFrame()
			if err != nil {
				done <- err
				return
			}

			ev := getEventByFrame(f)
			conn.vlog(ev, false)

			sf, ok := f.(*http2.SettingsFrame)
			if !ok {
				continue
			}

			if sf.IsAck() {
				local = true
			} else {
				remote = true
				sf.ForeachSetting(func(setting http2.Setting) error {
					conn.Settings[setting.ID] = setting.Val
					return nil
				})
				conn.WriteSettingsAck()
			}
		}

		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			return err
		}
	case <-time.After(conn.Timeout):
		return ErrTimeout
	}

	return nil
}

func (conn *Conn) MaxFrameSize() int {
	val, ok := conn.Settings[http2.SettingMaxFrameSize]
	if !ok {
		val = 16384
	}
	return int(val)
}

// EncodeHeaders encodes header and returns encoded bytes. Conn
// retains encoding context and next call of EncodeHeaders will be
// performed using the same encoding context.
func (conn *Conn) EncodeHeaders(headers []hpack.HeaderField) []byte {
	conn.encoderBuf.Reset()

	for _, hf := range headers {
		conn.encoder.WriteField(hf)
	}

	dst := make([]byte, conn.encoderBuf.Len())
	copy(dst, conn.encoderBuf.Bytes())

	return dst
}

func (conn *Conn) Send(payload []byte) error {
	conn.vlog(EventRawData{payload}, true)
	_, err := conn.Write(payload)
	return err
}

func (conn *Conn) WriteData(streamID uint32, endStream bool, data []byte) error {
	if conn.Verbose {
		conn.debugFramer.WriteData(streamID, endStream, data)
		conn.logFrameSend()
	}

	return conn.framer.WriteData(streamID, endStream, data)
}

func (conn *Conn) WriteHeaders(p http2.HeadersFrameParam) error {
	if conn.Verbose {
		conn.debugFramer.WriteHeaders(p)
		conn.logFrameSend()
	}

	return conn.framer.WriteHeaders(p)
}

func (conn *Conn) WritePriority(streamID uint32, p http2.PriorityParam) error {
	if conn.Verbose {
		conn.debugFramer.WritePriority(streamID, p)
		conn.logFrameSend()
	}

	return conn.framer.WritePriority(streamID, p)
}

func (conn *Conn) WriteRSTStream(streamID uint32, code http2.ErrCode) error {
	if conn.Verbose {
		conn.debugFramer.WriteRSTStream(streamID, code)
		conn.logFrameSend()
	}

	return conn.framer.WriteRSTStream(streamID, code)
}

func (conn *Conn) WriteSettings(settings ...http2.Setting) error {
	if conn.Verbose {
		conn.debugFramer.WriteSettings(settings...)
		conn.logFrameSend()
	}

	return conn.framer.WriteSettings(settings...)
}

func (conn *Conn) WriteSettingsAck() error {
	if conn.Verbose {
		conn.debugFramer.WriteSettingsAck()
		conn.logFrameSend()
	}

	return conn.framer.WriteSettingsAck()
}

func (conn *Conn) WritePing(ack bool, data [8]byte) error {
	if conn.Verbose {
		conn.debugFramer.WritePing(ack, data)
		conn.logFrameSend()
	}

	return conn.framer.WritePing(ack, data)
}

func (conn *Conn) WriteWindowUpdate(streamID, incr uint32) error {
	if conn.Verbose {
		conn.debugFramer.WriteWindowUpdate(streamID, incr)
		conn.logFrameSend()
	}

	return conn.framer.WriteWindowUpdate(streamID, incr)
}

func (conn *Conn) WriteContinuation(streamID uint32, endHeaders bool, headerBlockFragment []byte) error {
	if conn.Verbose {
		conn.debugFramer.WriteContinuation(streamID, endHeaders, headerBlockFragment)
		conn.logFrameSend()
	}

	return conn.framer.WriteContinuation(streamID, endHeaders, headerBlockFragment)
}

func (conn *Conn) WaitEvent() Event {
	var ev Event

	rd := time.Now().Add(conn.Timeout)
	conn.SetReadDeadline(rd)

	f, err := conn.framer.ReadFrame()
	if err != nil {
		if err == io.EOF {
			ev = EventConnectionClosed{}
			conn.vlog(ev, false)
			conn.Closed = true
			return ev
		}

		opErr, ok := err.(*net.OpError)
		if ok {
			if opErr.Err == syscall.ECONNRESET {
				ev = EventConnectionClosed{}
				conn.vlog(ev, false)
				conn.Closed = true
				return ev
			}

			if opErr.Timeout() {
				ev = EventTimeout{}
				conn.vlog(ev, false)
				conn.Closed = true
				return ev
			}
		}

		ev = EventError{}
		conn.vlog(ev, false)
		return ev
	}

	_, ok := f.(*http2.DataFrame)
	if ok {
		conn.updateWindowSize(f)
	}

	ev = getEventByFrame(f)
	conn.vlog(ev, false)

	return ev
}

func (conn *Conn) updateWindowSize(f http2.Frame) {
	if !conn.WindowUpdate {
		return
	}

	len := int(f.Header().Length)
	streamID := f.Header().StreamID

	_, ok := conn.WindowSize[streamID]
	if !ok {
		conn.WindowSize[streamID] = DefaultWindowSize
	}

	conn.WindowSize[streamID] -= len
	if conn.WindowSize[streamID] <= 0 {
		incr := DefaultWindowSize + (conn.WindowSize[streamID] * -1)
		conn.WriteWindowUpdate(streamID, uint32(incr))
		conn.WindowSize[streamID] += incr
	}

	conn.WindowSize[0] -= len
	if conn.WindowSize[0] <= 0 {
		incr := DefaultWindowSize + (conn.WindowSize[0] * -1)
		conn.WriteWindowUpdate(0, uint32(incr))
		conn.WindowSize[0] += incr
	}
}

func (conn *Conn) logFrameSend() {
	f, err := conn.debugFramer.ReadFrame()
	if err != nil {
		// http2 package does not parse DATA frame with stream ID: 0x0.
		// So we are going to log the information that sent some frame.
		if conn.Verbose {
			log.Println(gray(fmt.Sprintf("     <-- [send] ??? Frame (Failed to parse the frame)")))
		}
		return
	}

	ev := getEventByFrame(f)
	conn.vlog(ev, true)
}

func (conn *Conn) vlog(ev Event, send bool) {
	if !conn.Verbose {
		return
	}

	if send {
		log.Println(gray(fmt.Sprintf("     <-- [send] %s", ev)))
	} else {
		log.Println(gray(fmt.Sprintf("     --> [recv] %s", ev)))
	}
}

func getEventByFrame(f http2.Frame) Event {
	var ev Event

	switch f := f.(type) {
	case *http2.DataFrame:
		ev = EventDataFrame{*f}
	case *http2.HeadersFrame:
		ev = EventHeadersFrame{*f}
	case *http2.PriorityFrame:
		ev = EventPriorityFrame{*f}
	case *http2.RSTStreamFrame:
		ev = EventRSTStreamFrame{*f}
	case *http2.SettingsFrame:
		ev = EventSettingsFrame{*f}
	case *http2.PushPromiseFrame:
		ev = EventPushPromiseFrame{*f}
	case *http2.PingFrame:
		ev = EventPingFrame{*f}
	case *http2.GoAwayFrame:
		ev = EventGoAwayFrame{*f}
	case *http2.WindowUpdateFrame:
		ev = EventWindowUpdateFrame{*f}
	case *http2.ContinuationFrame:
		ev = EventContinuationFrame{*f}
		//default:
		//	ev = EventUnknownFrame(f)
	}

	return ev
}
