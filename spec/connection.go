package spec

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
)

const (
	// DefaultWindowSize is the value of default connection window size.
	DefaultWindowSize = 65535
	// DefaultFrameSize is the value of default frame size.
	DefaultFrameSize = 16384
)

// Conn represent a HTTP/2 connection.
// This struct contains settings information, current window size,
// encoder of HPACK and frame encoder.
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
	decoder    *hpack.Decoder

	debugFramer    *http2.Framer
	debugFramerBuf *bytes.Buffer

	server bool
}

// Dial connects to the server based on configuration.
func Dial(c *config.Config) (*Conn, error) {
	var baseConn net.Conn
	var err error

	if c.TLS {
		dialer := &net.Dialer{}
		dialer.Timeout = c.Timeout

		tlsConfig, err := c.TLSConfig()
		if err != nil {
			return nil, err
		}

		tlsConn, err := tls.DialWithDialer(dialer, "tcp", c.Addr(), tlsConfig)
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

	decoder := hpack.NewDecoder(4096, func(f hpack.HeaderField) {})

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
		decoder:    decoder,

		server: false,
	}

	if conn.Verbose {
		conn.debugFramerBuf = new(bytes.Buffer)
		conn.debugFramer = http2.NewFramer(conn.debugFramerBuf, conn.debugFramerBuf)
		conn.debugFramer.AllowIllegalWrites = true
		conn.debugFramer.AllowIllegalReads = true
	}

	return &conn, nil
}

func Accept(c *config.Config, baseConn net.Conn) (*Conn, error) {
	settings := map[http2.SettingID]uint32{}

	framer := http2.NewFramer(baseConn, baseConn)
	framer.AllowIllegalWrites = true
	framer.AllowIllegalReads = true

	var encoderBuf bytes.Buffer
	encoder := hpack.NewEncoder(&encoderBuf)

	decoder := hpack.NewDecoder(4096, func(f hpack.HeaderField) {})

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
		decoder:    decoder,

		server: true,
	}

	if conn.Verbose {
		conn.debugFramerBuf = new(bytes.Buffer)
		conn.debugFramer = http2.NewFramer(conn.debugFramerBuf, conn.debugFramerBuf)
		conn.debugFramer.AllowIllegalWrites = true
		conn.debugFramer.AllowIllegalReads = true
	}

	return &conn, nil
}

// Handshake performs HTTP/2 handshake with the server.
func (conn *Conn) Handshake() error {
	if conn.server {
		return conn.handshakeAsServer()
	} else {
		return conn.handshakeAsClient()
	}
}

// MaxFrameSize returns value of Handshake performs HTTP/2 handshake
// with the server.
func (conn *Conn) MaxFrameSize() int {
	val, ok := conn.Settings[http2.SettingMaxFrameSize]
	if !ok {
		return DefaultFrameSize
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

// SetMaxDynamicTableSize changes the dynamic header table size to v.
func (conn *Conn) SetMaxDynamicTableSize(v uint32) {
	conn.encoder.SetMaxDynamicTableSize(v)
}

// Send sends a byte sequense. This function is used to send a raw
// data in tests.
func (conn *Conn) Send(payload []byte) error {
	conn.vlog(RawDataEvent{payload}, true)
	_, err := conn.Write(payload)
	return err
}

// WriteData sends a DATA frame.
func (conn *Conn) WriteData(streamID uint32, endStream bool, data []byte) error {
	if conn.Verbose {
		conn.debugFramer.WriteData(streamID, endStream, data)
		conn.logFrameSend()
	}

	return conn.framer.WriteData(streamID, endStream, data)
}

// WriteDataPadded sends a DATA frame with padding.
func (conn *Conn) WriteDataPadded(streamID uint32, endStream bool, data, pad []byte) error {
	if conn.Verbose {
		conn.debugFramer.WriteDataPadded(streamID, endStream, data, pad)
		conn.logFrameSend()
	}

	return conn.framer.WriteDataPadded(streamID, endStream, data, pad)
}

// WriteHeaders sends a HEADERS frame.
func (conn *Conn) WriteHeaders(p http2.HeadersFrameParam) error {
	if conn.Verbose {
		conn.debugFramer.WriteHeaders(p)
		conn.logFrameSend()
	}

	return conn.framer.WriteHeaders(p)
}

// WritePriority sends a PRIORITY frame.
func (conn *Conn) WritePriority(streamID uint32, p http2.PriorityParam) error {
	if conn.Verbose {
		conn.debugFramer.WritePriority(streamID, p)
		conn.logFrameSend()
	}

	return conn.framer.WritePriority(streamID, p)
}

// WriteRSTStream sends a RST_STREAM frame.
func (conn *Conn) WriteRSTStream(streamID uint32, code http2.ErrCode) error {
	if conn.Verbose {
		conn.debugFramer.WriteRSTStream(streamID, code)
		conn.logFrameSend()
	}

	return conn.framer.WriteRSTStream(streamID, code)
}

// WriteSettings sends a SETTINGS frame.
func (conn *Conn) WriteSettings(settings ...http2.Setting) error {
	if conn.Verbose {
		conn.debugFramer.WriteSettings(settings...)
		conn.logFrameSend()
	}

	return conn.framer.WriteSettings(settings...)
}

// WriteSettingsAck sends a SETTINGS frame with ACK flag.
func (conn *Conn) WriteSettingsAck() error {
	if conn.Verbose {
		conn.debugFramer.WriteSettingsAck()
		conn.logFrameSend()
	}

	return conn.framer.WriteSettingsAck()
}

// WritePushPromise sends a PUSH_PROMISE frame.
func (conn *Conn) WritePushPromise(p http2.PushPromiseParam) error {
	if conn.Verbose {
		conn.debugFramer.WritePushPromise(p)
		conn.logFrameSend()
	}

	return conn.framer.WritePushPromise(p)
}

// WritePing sends a PING frame.
func (conn *Conn) WritePing(ack bool, data [8]byte) error {
	if conn.Verbose {
		conn.debugFramer.WritePing(ack, data)
		conn.logFrameSend()
	}

	return conn.framer.WritePing(ack, data)
}

// WriteGoAway sends a GOAWAY frame.
func (conn *Conn) WriteGoAway(maxStreamID uint32, code http2.ErrCode, debugData []byte) error {
	if conn.Verbose {
		conn.debugFramer.WriteGoAway(maxStreamID, code, debugData)
		conn.logFrameSend()
	}

	return conn.framer.WriteGoAway(maxStreamID, code, debugData)
}

// WriteWindowUpdate sends a WINDOW_UPDATE frame.
func (conn *Conn) WriteWindowUpdate(streamID, incr uint32) error {
	if conn.Verbose {
		conn.debugFramer.WriteWindowUpdate(streamID, incr)
		conn.logFrameSend()
	}

	return conn.framer.WriteWindowUpdate(streamID, incr)
}

// WriteContinuation sends a CONTINUATION frame.
func (conn *Conn) WriteContinuation(streamID uint32, endHeaders bool, headerBlockFragment []byte) error {
	if conn.Verbose {
		conn.debugFramer.WriteContinuation(streamID, endHeaders, headerBlockFragment)
		conn.logFrameSend()
	}

	return conn.framer.WriteContinuation(streamID, endHeaders, headerBlockFragment)
}

func (conn *Conn) WriteRawFrame(t http2.FrameType, flags http2.Flags, streamID uint32, payload []byte) error {
	if conn.Verbose {
		conn.debugFramer.WriteRawFrame(t, flags, streamID, payload)
		conn.logFrameSend()
	}

	return conn.framer.WriteRawFrame(t, flags, streamID, payload)
}

func (conn *Conn) WriteSuccessResponse(streamID uint32, c *config.Config) {
	hp := http2.HeadersFrameParam{
		StreamID:      streamID,
		EndStream:     false,
		EndHeaders:    true,
		BlockFragment: conn.EncodeHeaders(CommonRespHeaders(c)),
	}
	conn.WriteHeaders(hp)
	conn.WriteData(streamID, true, []byte("success"))
}

// WaitEvent returns a event occured on connection. This function is
// used to wait the next event on the connection.
func (conn *Conn) WaitEvent() Event {
	var ev Event

	rd := time.Now().Add(conn.Timeout)
	conn.SetReadDeadline(rd)

	f, err := conn.framer.ReadFrame()
	if err != nil {
		conn.Closed = true

		if err == io.EOF {
			ev = ConnectionClosedEvent{}
			conn.vlog(ev, false)
			return ev
		}

		opErr, ok := err.(*net.OpError)
		if ok {
			if opErr.Err == syscall.ECONNRESET {
				ev = ConnectionClosedEvent{}
				conn.vlog(ev, false)
				return ev
			}

			if runtime.GOOS == "windows" {
				scErr, ok := opErr.Err.(*os.SyscallError)
				if ok {
					const WSAECONNABORTED = 10053
					const WSAECONNRESET = 10054

					rv := reflect.ValueOf(scErr.Err)
					if rv.Kind() == reflect.Uintptr {
						n := uintptr(rv.Uint())
						if n == WSAECONNRESET || n == WSAECONNABORTED {
							ev = ConnectionClosedEvent{}
							conn.vlog(ev, false)
							return ev
						}
					}
				}
			}

			if opErr.Timeout() {
				ev = TimeoutEvent{}
				conn.vlog(ev, false)
				return ev
			}
		}

		ev = ErrorEvent{err}
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

// WaitEventByType returns a specified event occured on connection.
// This function is used to wait the next event that has specified
// type on the connection.
func (conn *Conn) WaitEventByType(evt EventType) (Event, bool) {
	var lastEvent Event

	for !conn.Closed {
		ev := conn.WaitEvent()

		if ev.Type() == evt {
			return ev, true
		}

		if ev.Type() == EventTimeout && lastEvent != nil {
			break
		}

		lastEvent = ev
	}

	return lastEvent, false
}

type Request struct {
	StreamID uint32
	Headers  []hpack.HeaderField
}

func (conn *Conn) ReadRequest() (*Request, error) {
	headers := make([]hpack.HeaderField, 0, 256)
	conn.decoder.SetEmitFunc(func(f hpack.HeaderField) {
		headers = append(headers, f)
	})

	done := false
	streamID := uint32(0)

	for !done {
		f, ok := conn.WaitEventByType(EventHeadersFrame)
		if !ok {
			return nil, errors.New("No HEADER frame received")
		}

		hf, _ := f.(HeadersFrameEvent)

		if streamID == uint32(0) {
			streamID = hf.Header().StreamID
		} else if streamID != hf.Header().StreamID {
			return nil, errors.New("Encountered different StreamID")
		}

		_, err := conn.decoder.Write(hf.HeaderBlockFragment())
		if err != nil {
			return nil, err
		}

		done = hf.HeadersEnded()
	}

	request := &Request{
		StreamID: streamID,
		Headers:  headers,
	}
	return request, nil
}

// updateWindowSize calculates the current window size based on the
// information in the given HTTP/2 frame.
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

// logFrameSend writes a log of the frame to be sent.
func (conn *Conn) logFrameSend() {
	f, err := conn.debugFramer.ReadFrame()
	if err != nil {
		// http2 package does not parse DATA frame with stream ID: 0x0.
		// So we are going to log the information that sent some frame.
		if conn.Verbose {
			log.Println(gray(fmt.Sprintf("     [send] ??? Frame (Failed to parse the frame)")))
		}
		return
	}

	ev := getEventByFrame(f)
	conn.vlog(ev, true)
}

// vlog writes a verbose log.
func (conn *Conn) vlog(ev Event, send bool) {
	if !conn.Verbose {
		return
	}

	if send {
		log.Println(gray(fmt.Sprintf("     [send] %s", ev)))
	} else {
		log.Println(gray(fmt.Sprintf("     [recv] %s", ev)))
	}
}

// getEventByFrame returns an event based on given HTTP/2 frame.
func getEventByFrame(f http2.Frame) Event {
	var ev Event

	switch f := f.(type) {
	case *http2.DataFrame:
		ev = DataFrameEvent{*f}
	case *http2.HeadersFrame:
		ev = HeadersFrameEvent{*f}
	case *http2.PriorityFrame:
		ev = PriorityFrameEvent{*f}
	case *http2.RSTStreamFrame:
		ev = RSTStreamFrameEvent{*f}
	case *http2.SettingsFrame:
		ev = SettingsFrameEvent{*f}
	case *http2.PushPromiseFrame:
		ev = PushPromiseFrameEvent{*f}
	case *http2.PingFrame:
		ev = PingFrameEvent{*f}
	case *http2.GoAwayFrame:
		ev = GoAwayFrameEvent{*f}
	case *http2.WindowUpdateFrame:
		ev = WindowUpdateFrameEvent{*f}
	case *http2.ContinuationFrame:
		ev = ContinuationFrameEvent{*f}
		//default:
		//	ev = EventUnknownFrame(f)
	}

	return ev
}

func (conn *Conn) handshakeAsClient() error {
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

func (conn *Conn) handshakeAsServer() error {
	done := make(chan error)

	go func() {
		_, err := conn.ReadClientPreface()
		if err != nil {
			done <- err
			return
		}

		f, err := conn.framer.ReadFrame()
		if err != nil {
			done <- err
			return
		}

		_, ok := f.(*http2.SettingsFrame)
		if !ok {
			done <- errors.New("First frame must be SETTINGS frame")
			return
		}

		setting := http2.Setting{
			ID:  http2.SettingInitialWindowSize,
			Val: DefaultWindowSize,
		}
		conn.WriteSettings(setting)
		conn.WriteSettingsAck()

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

func (conn *Conn) ReadClientPreface() (string, error) {
	prefaceBytes, err := conn.readBytes(24)
	if err != nil {
		return "", err
	}

	preface := string(prefaceBytes[:])
	if preface != "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n" {
		return "", errors.New("Illegal preface")
	}
	return preface, nil
}

func (conn *Conn) readBytes(size int) ([]byte, error) {
	var remain = size
	buffer := make([]byte, 0, size)

	for remain > 0 {
		tmp := make([]byte, remain)
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, err
		}

		buffer = append(buffer, tmp[:n]...)
		remain = remain - n
	}
	return buffer, nil
}
