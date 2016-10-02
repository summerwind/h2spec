package spec

import (
	"math"

	"golang.org/x/net/http2"
)

var (
	DefaultLength  uint32        = math.MaxUint32
	DefaultFlags   http2.Flags   = math.MaxUint8
	DefaultErrCode http2.ErrCode = math.MaxUint8
)

type Event interface {
	String() string
}

type EventConnectionClosed struct{}

func (ev EventConnectionClosed) String() string {
	return "Connection closed"
}

type EventError struct{}

func (ev EventError) String() string {
	return "Error"
}

type EventTimeout struct {
	LastEvent Event
}

func (ev EventTimeout) String() string {
	return "Timeout"
}

type EventDataFrame struct {
	http2.DataFrame
}

func (ev EventDataFrame) String() string {
	return "DATA Frame"
}

type EventHeadersFrame struct {
	http2.HeadersFrame
}

func (ev EventHeadersFrame) String() string {
	return "HEADERS Frame"
}

type EventPriorityFrame struct {
	http2.PriorityFrame
}

func (ev EventPriorityFrame) String() string {
	return "PRIORITY Frame"
}

type EventRSTStreamFrame struct {
	http2.RSTStreamFrame
}

func (ev EventRSTStreamFrame) String() string {
	return "RST_STREAM Frame"
}

type EventSettingsFrame struct {
	http2.SettingsFrame
}

func (ev EventSettingsFrame) String() string {
	return "SETTINGS Frame"
}

type EventPushPromiseFrame struct {
	http2.PushPromiseFrame
}

func (ev EventPushPromiseFrame) String() string {
	return "PUSH_PROMISE Frame"
}

type EventPingFrame struct {
	http2.PingFrame
}

func (ev EventPingFrame) String() string {
	return "PING Frame"
}

type EventGoAwayFrame struct {
	http2.GoAwayFrame
}

func (ev EventGoAwayFrame) String() string {
	return "GOAWAY Frame"
}

type EventWindowUpdateFrame struct {
	http2.WindowUpdateFrame
}

func (ev EventWindowUpdateFrame) String() string {
	return "WINDOW_UPDATE Frame"
}

type EventContinuationFrame struct {
	http2.ContinuationFrame
}

func (ev EventContinuationFrame) String() string {
	return "CONTINUATION Frame"
}
