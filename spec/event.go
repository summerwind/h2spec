package spec

import (
	"fmt"
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

type EventFrame interface {
	String() string
	Header() http2.FrameHeader
}

type EventConnectionClosed struct{}

func (ev EventConnectionClosed) String() string {
	return "Connection closed"
}

type EventError struct {
	Error error
}

func (ev EventError) String() string {
	return fmt.Sprintf("Error: %v", ev.Error)
}

type EventTimeout struct {
	LastEvent Event
}

func (ev EventTimeout) String() string {
	return "Timeout"
}

type EventRawData struct {
	Payload []byte
}

func (ev EventRawData) String() string {
	return fmt.Sprintf("Raw Data (0x%x)", ev.Payload)
}

type EventDataFrame struct {
	http2.DataFrame
}

func (ev EventDataFrame) String() string {
	return frameString(ev.Header())
}

type EventHeadersFrame struct {
	http2.HeadersFrame
}

func (ev EventHeadersFrame) String() string {
	return frameString(ev.Header())
}

type EventPriorityFrame struct {
	http2.PriorityFrame
}

func (ev EventPriorityFrame) String() string {
	return frameString(ev.Header())
}

type EventRSTStreamFrame struct {
	http2.RSTStreamFrame
}

func (ev EventRSTStreamFrame) String() string {
	return frameString(ev.Header())
}

type EventSettingsFrame struct {
	http2.SettingsFrame
}

func (ev EventSettingsFrame) String() string {
	return frameString(ev.Header())
}

type EventPushPromiseFrame struct {
	http2.PushPromiseFrame
}

func (ev EventPushPromiseFrame) String() string {
	return frameString(ev.Header())
}

type EventPingFrame struct {
	http2.PingFrame
}

func (ev EventPingFrame) String() string {
	return frameString(ev.Header())
}

type EventGoAwayFrame struct {
	http2.GoAwayFrame
}

func (ev EventGoAwayFrame) String() string {
	return frameString(ev.Header())
}

type EventWindowUpdateFrame struct {
	http2.WindowUpdateFrame
}

func (ev EventWindowUpdateFrame) String() string {
	return frameString(ev.Header())
}

type EventContinuationFrame struct {
	http2.ContinuationFrame
}

func (ev EventContinuationFrame) String() string {
	return frameString(ev.Header())
}

func frameString(header http2.FrameHeader) string {
	return fmt.Sprintf(
		"%s Frame (length:%d, flags:0x%02x, stream_id:%d)",
		header.Type,
		header.Length,
		header.Flags,
		header.StreamID,
	)
}
