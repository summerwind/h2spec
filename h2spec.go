package h2spec

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/bradfitz/http2"
	"github.com/bradfitz/http2/hpack"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var TIMEOUT = errors.New("Timeout")

type RunMode int

const (
	ModeAll       RunMode = 0
	ModeGroupOnly RunMode = 1
	ModeSkip      RunMode = 2
)

type Context struct {
	Port      int
	Host      string
	Tls       bool
	TlsConfig *tls.Config
	Sections  map[string]bool
	Timeout   time.Duration
}

func (ctx *Context) Authority() string {
	return fmt.Sprintf("%s:%d", ctx.Host, ctx.Port)
}

func (ctx *Context) GetRunMode(section string) RunMode {
	if ctx.Sections == nil {
		return ModeAll
	}

	val, ok := ctx.Sections[section]
	if !ok {
		return ModeSkip
	}
	if !val {
		return ModeGroupOnly
	}

	return ModeAll
}

type Test interface {
	Run(*Context, int)
}

type TestGroup struct {
	Section    string
	Name       string
	testGroups []TestGroup
	testCases  []TestCase
}

func (tg *TestGroup) Run(ctx *Context, level int) {
	runMode := ctx.GetRunMode(tg.Section)

	if runMode == ModeSkip {
		return
	}

	tg.PrintHeader(level)
	if runMode == ModeAll {
		for _, testCase := range tg.testCases {
			testCase.Run(ctx, level+1)
		}
		tg.PrintFooter(level)
	}

	for _, testGroup := range tg.testGroups {
		testGroup.Run(ctx, level+1)
	}
}

func (tg *TestGroup) AddTestCase(testCase *TestCase) {
	tg.testCases = append(tg.testCases, *testCase)
}

func (tg *TestGroup) AddTestGroup(testGroup *TestGroup) {
	tg.testGroups = append(tg.testGroups, *testGroup)
}

func (tg *TestGroup) PrintHeader(level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s%s. %s\n", indent, tg.Section, tg.Name)
}

func (tg *TestGroup) PrintFooter(level int) {
	if len(tg.testCases) == 0 {
		return
	}
	fmt.Println("")
}

type TestCase struct {
	Desc    string
	Spec    string
	handler func(*Context) ([]Result, Result)
}

func (tc *TestCase) Run(ctx *Context, level int) {
	expected, actual := tc.handler(ctx)

	_, ok := actual.(*ResultSkipped)
	if ok {
		tc.PrintSkipped(actual, level)
		return
	}

	if tc.evaluateResult(expected, actual) {
		tc.PrintResult(level)
	} else {
		tc.PrintError(expected, actual, level)
	}
}

func (tc *TestCase) evaluateResult(expected []Result, actual Result) bool {
	actualStr := actual.String()
	for _, exp := range expected {
		if exp.String() == actualStr {
			return true
		}
	}

	return false
}

func (tc *TestCase) HandleFunc(handler func(*Context) ([]Result, Result)) {
	tc.handler = handler
}

func (tc *TestCase) PrintResult(level int) {
	mark := "✓"
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s\x1b[32m%s\x1b[0m \x1b[90m%s\x1b[0m\n", indent, mark, tc.Desc)
}

func (tc *TestCase) PrintError(expected []Result, actual Result, level int) {
	mark := "×"
	indent := strings.Repeat("  ", level)

	fmt.Printf("\x1b[31m")
	fmt.Printf("%s%s %s\n", indent, mark, tc.Desc)
	fmt.Printf("%s  - %s\n", indent, tc.Spec)
	fmt.Printf("\x1b[32m")
	for i, exp := range expected {
		var lavel string
		if i == 0 {
			lavel = "Expected:"
		} else {
			lavel = strings.Repeat(" ", 9)
		}
		fmt.Printf("%s    %s %s\n", indent, lavel, exp)
	}
	fmt.Printf("\x1b[33m")
	fmt.Printf("%s      Actual: %s\n", indent, actual)
	fmt.Printf("\x1b[0m")
}

func (tc *TestCase) PrintSkipped(actual Result, level int) {
	mark := " "
	indent := strings.Repeat("  ", level)

	fmt.Printf("\x1b[36m")
	fmt.Printf("%s%s %s\n", indent, mark, tc.Desc)
	fmt.Printf("%s  - %s\n", indent, actual)
	fmt.Printf("\x1b[0m")
}

func NewTestGroup(section, name string) *TestGroup {
	return &TestGroup{section, name, nil, nil}
}

func NewTestCase(desc, spec string, handler func(*Context) ([]Result, Result)) *TestCase {
	return &TestCase{desc, spec, handler}
}

var FlagDefault http2.Flags = 0x0
var ErrCodeDefault http2.ErrCode = 0xff

type Result interface {
	String() string
}

type ResultFrame struct {
	Type  http2.FrameType
	Flags http2.Flags
	Code  http2.ErrCode
}

func (rf *ResultFrame) String() string {
	parts := []string{}

	if rf.Flags != FlagDefault {
		parts = append(parts, fmt.Sprintf("Flags: %d", rf.Flags))
	}
	if rf.Code != ErrCodeDefault {
		parts = append(parts, fmt.Sprintf("ErrorCode: %s", rf.Code.String()))
	}

	res := fmt.Sprintf("%s frame", rf.Type.String())
	if len(parts) > 0 {
		res += fmt.Sprintf(" (%s)", strings.Join(parts, ", "))
	}

	return res
}

type ResultConnectionClose struct{}

func (rcc *ResultConnectionClose) String() string {
	return "Connection close"
}

type ResultTestTimeout struct{}

func (ttr *ResultTestTimeout) String() string {
	return "Test timeout"
}

type ResultSkipped struct {
	Reason string
}

func (rs *ResultSkipped) String() string {
	return "Skipped: " + rs.Reason
}

type ResultError struct {
	Error error
}

func (re *ResultError) String() string {
	return fmt.Sprintf("Error: %s", re.Error)
}

type TcpConn struct {
	conn   net.Conn
	dataCh chan []byte
	errCh  chan error
}

type Http2Conn struct {
	conn           net.Conn
	dataCh         chan http2.Frame
	errCh          chan error
	fr             *http2.Framer
	HpackEncoder   *hpack.Encoder
	HeaderWriteBuf bytes.Buffer
	Settings       map[http2.SettingID]uint32
}

// ReadFrame reads a complete HTTP/2 frame from underlying connection.
// This function blocks until a complete frame is received or timeout
// t is expired.  The returned http2.Frame must not be used after next
// ReadFrame call.
func (h2Conn *Http2Conn) ReadFrame(t time.Duration) (http2.Frame, error) {
	go func() {
		f, err := h2Conn.fr.ReadFrame()
		if err != nil {
			h2Conn.errCh <- err
			return
		}
		h2Conn.dataCh <- f
	}()

	select {
	case f := <-h2Conn.dataCh:
		return f, nil
	case err := <-h2Conn.errCh:
		return nil, err
	case <-time.After(t):
		return nil, TIMEOUT
	}
}

// EncodeHeader encodes header and returns encoded bytes.  h2Conn
// retains encoding context and next call of EncodeHeader will be
// performed using the same encoding context.
func (h2Conn *Http2Conn) EncodeHeader(header []hpack.HeaderField) []byte {
	h2Conn.HeaderWriteBuf.Reset()

	for _, hf := range header {
		_ = h2Conn.HpackEncoder.WriteField(hf)
	}

	dst := make([]byte, h2Conn.HeaderWriteBuf.Len())
	copy(dst, h2Conn.HeaderWriteBuf.Bytes())

	return dst
}

func connectTls(ctx *Context) (net.Conn, error) {
	if ctx.TlsConfig == nil {
		ctx.TlsConfig = new(tls.Config)
	}

	if ctx.TlsConfig.NextProtos == nil {
		ctx.TlsConfig.NextProtos = append(ctx.TlsConfig.NextProtos, "h2-14", "h2-15", "h2-16")
	}

	conn, err := tls.Dial("tcp", ctx.Authority(), ctx.TlsConfig)
	if err != nil {
		return nil, err
	}

	cs := conn.ConnectionState()
	if !cs.NegotiatedProtocolIsMutual {
		return nil, fmt.Errorf("HTTP/2 protocol was not negotiated")
	}

	return conn, err
}

func CreateTcpConn(ctx *Context) *TcpConn {
	var conn net.Conn
	var err error

	if ctx.Tls {
		conn, err = connectTls(ctx)
	} else {
		conn, err = net.Dial("tcp", ctx.Authority())
	}

	if err != nil {
		fmt.Printf("Unable to connect to the target server: %v\n", err)
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
	var conn net.Conn
	var err error

	if ctx.Tls {
		conn, err = connectTls(ctx)
	} else {
		conn, err = net.Dial("tcp", ctx.Authority())
	}

	if err != nil {
		fmt.Printf("Unable to connect to the target server: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(conn, "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")

	fr := http2.NewFramer(conn, conn)
	settings := map[http2.SettingID]uint32{}

	if sn {
		doneCh := make(chan bool, 1)
		errCh := make(chan error, 1)
		fr.WriteSettings()

		go func() {
			local := false
			remote := false

			for {
				f, err := fr.ReadFrame()
				if err != nil {
					errCh <- err
					return
				}

				switch f := f.(type) {
				case *http2.SettingsFrame:
					if f.IsAck() {
						local = true
					} else {
						f.ForeachSetting(func(setting http2.Setting) error {
							settings[setting.ID] = setting.Val
							return nil
						})
						fr.WriteSettingsAck()
						remote = true
					}
				}

				if local && remote {
					doneCh <- true
					return
				}
			}
		}()

		select {
		case <-doneCh:
			// Nothing to. do
		case <-errCh:
			fmt.Println("HTTP/2 settings negotiation failed")
			os.Exit(1)
		case <-time.After(ctx.Timeout):
			fmt.Println("HTTP/2 settings negotiation timeout")
			os.Exit(1)
		}
	}

	fr.AllowIllegalWrites = true
	dataCh := make(chan http2.Frame)
	errCh := make(chan error, 1)

	http2Conn := &Http2Conn{
		conn:     conn,
		fr:       fr,
		dataCh:   dataCh,
		errCh:    errCh,
		Settings: settings,
	}

	http2Conn.HpackEncoder = hpack.NewEncoder(&http2Conn.HeaderWriteBuf)

	return http2Conn
}

func TestConnectionError(ctx *Context, http2Conn *Http2Conn, codes []http2.ErrCode) (expected []Result, actual Result) {
	for _, code := range codes {
		expected = append(expected, &ResultFrame{http2.FrameGoAway, FlagDefault, code})
	}
	expected = append(expected, &ResultConnectionClose{})

loop:
	for {
		f, err := http2Conn.ReadFrame(ctx.Timeout)
		if err != nil {
			if err == io.EOF {
				actual = &ResultConnectionClose{}
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
			actual = &ResultFrame{f.Header().Type, FlagDefault, f.ErrCode}
			if TestErrorCode(f.ErrCode, codes) {
				break loop
			}
		default:
			actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
		}
	}

	return expected, actual
}

func TestStreamError(ctx *Context, http2Conn *Http2Conn, codes []http2.ErrCode) (expected []Result, actual Result) {
	for _, code := range codes {
		expected = append(expected, &ResultFrame{http2.FrameGoAway, FlagDefault, code})
		expected = append(expected, &ResultFrame{http2.FrameRSTStream, FlagDefault, code})
	}
	expected = append(expected, &ResultConnectionClose{})

loop:
	for {
		f, err := http2Conn.ReadFrame(ctx.Timeout)
		if err != nil {
			if err == io.EOF {
				actual = &ResultConnectionClose{}
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
			actual = &ResultFrame{f.Header().Type, FlagDefault, f.ErrCode}
			if TestErrorCode(f.ErrCode, codes) {
				break loop
			}
		case *http2.RSTStreamFrame:
			actual = &ResultFrame{f.Header().Type, FlagDefault, f.ErrCode}
			if TestErrorCode(f.ErrCode, codes) {
				break loop
			}
		default:
			actual = &ResultFrame{f.Header().Type, FlagDefault, ErrCodeDefault}
		}
	}

	return expected, actual
}

func TestErrorCode(code http2.ErrCode, expected []http2.ErrCode) bool {
	for _, exp := range expected {
		if code == exp {
			return true
		}
	}
	return false
}

func dummyData(num int) string {
	var data string
	for i := 0; i < num; i++ {
		data += "x"
	}
	return data
}

func pair(name, value string) hpack.HeaderField {
	return hpack.HeaderField{Name: name, Value: value}
}

func Run(ctx *Context) {
	groups := []*TestGroup{
		Http2ConnectionPrefaceTestGroup(),
		FrameSizeTestGroup(),
		HeaderCompressionAndDecompressionTestGroup(),
		StreamStatesTestGroup(),
		ErrorHandlingTestGroup(),
		ExtendingHttp2TestGroup(),
		DataTestGroup(),
		HeadersTestGroup(),
		PriorityTestGroup(),
		RstStreamTestGroup(),
		SettingsTestGroup(),
		PingTestGroup(),
		GoawayTestGroup(),
		WindowUpdateTestGroup(),
		ContinuationTestGroup(),
		HttpRequestResponseExchangeTestGroup(),
		ServerPushTestGroup(),
	}

	for _, group := range groups {
		group.Run(ctx, 1)
	}
}
