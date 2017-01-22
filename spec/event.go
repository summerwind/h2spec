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

type EventType uint8

const (
	EventDataFrame         EventType = 0x0
	EventHeadersFrame      EventType = 0x1
	EventPriorityFrame     EventType = 0x2
	EventRSTStreamFrame    EventType = 0x3
	EventSettingsFrame     EventType = 0x4
	EventPushPromiseFrame  EventType = 0x5
	EventPingFrame         EventType = 0x6
	EventGoAwayFrame       EventType = 0x7
	EventWindowUpdateFrame EventType = 0x8
	EventContinuationFrame EventType = 0x9
	EventRawData           EventType = 0x10
	EventConnectionClosed  EventType = 0x11
	EventError             EventType = 0x12
	EventTimeout           EventType = 0x13
)

var eventName = map[EventType]string{
	EventDataFrame:         "DATA frame",
	EventHeadersFrame:      "HEADERS frame",
	EventPriorityFrame:     "PRIORITY frame",
	EventRSTStreamFrame:    "RST_STREAM frame",
	EventSettingsFrame:     "SETTINGS frame",
	EventPushPromiseFrame:  "PUSH_PROMISE frame",
	EventPingFrame:         "PING frame",
	EventGoAwayFrame:       "GOAWAY frame",
	EventWindowUpdateFrame: "WINDOW_UPDATE frame",
	EventContinuationFrame: "CONTINUATION frame",
	EventRawData:           "Raw data",
	EventConnectionClosed:  "Connection closed",
	EventError:             "Error",
	EventTimeout:           "Timeout",
}

func (et EventType) String() string {
	s, ok := eventName[et]
	if ok {
		return s
	}
	return fmt.Sprintf("Unknown event (%d)", uint8(et))
}

type Event interface {
	Type() EventType
	String() string
}

type EventFrame interface {
	String() string
	Header() http2.FrameHeader
}

type ConnectionClosedEvent struct{}

func (ev ConnectionClosedEvent) Type() EventType {
	return EventConnectionClosed
}

func (ev ConnectionClosedEvent) String() string {
	return "Connection closed"
}

type ErrorEvent struct {
	Error error
}

func (ev ErrorEvent) Type() EventType {
	return EventError
}

func (ev ErrorEvent) String() string {
	return fmt.Sprintf("Error: %v", ev.Error)
}

type TimeoutEvent struct{}

func (ev TimeoutEvent) Type() EventType {
	return EventTimeout
}

func (ev TimeoutEvent) String() string {
	return "Timeout"
}

type RawDataEvent struct {
	Payload []byte
}

func (ev RawDataEvent) Type() EventType {
	return EventRawData
}

func (ev RawDataEvent) String() string {
	return fmt.Sprintf("Raw Data (0x%x)", ev.Payload)
}

type DataFrameEvent struct {
	http2.DataFrame
}

func (ev DataFrameEvent) Type() EventType {
	return EventDataFrame
}

func (ev DataFrameEvent) String() string {
	return frameString(ev.Header())
}

type HeadersFrameEvent struct {
	http2.HeadersFrame
}

func (ev HeadersFrameEvent) Type() EventType {
	return EventHeadersFrame
}

func (ev HeadersFrameEvent) String() string {
	return frameString(ev.Header())
}

type PriorityFrameEvent struct {
	http2.PriorityFrame
}

func (ev PriorityFrameEvent) Type() EventType {
	return EventPriorityFrame
}

func (ev PriorityFrameEvent) String() string {
	return frameString(ev.Header())
}

type RSTStreamFrameEvent struct {
	http2.RSTStreamFrame
}

func (ev RSTStreamFrameEvent) Type() EventType {
	return EventRSTStreamFrame
}

func (ev RSTStreamFrameEvent) String() string {
	return frameString(ev.Header())
}

type SettingsFrameEvent struct {
	http2.SettingsFrame
}

func (ev SettingsFrameEvent) Type() EventType {
	return EventSettingsFrame
}

func (ev SettingsFrameEvent) String() string {
	return frameString(ev.Header())
}

type PushPromiseFrameEvent struct {
	http2.PushPromiseFrame
}

func (ev PushPromiseFrameEvent) Type() EventType {
	return EventPushPromiseFrame
}

func (ev PushPromiseFrameEvent) String() string {
	return frameString(ev.Header())
}

type PingFrameEvent struct {
	http2.PingFrame
}

func (ev PingFrameEvent) Type() EventType {
	return EventPingFrame
}

func (ev PingFrameEvent) String() string {
	return frameString(ev.Header())
}

type GoAwayFrameEvent struct {
	http2.GoAwayFrame
}

func (ev GoAwayFrameEvent) Type() EventType {
	return EventGoAwayFrame
}

func (ev GoAwayFrameEvent) String() string {
	return frameString(ev.Header())
}

type WindowUpdateFrameEvent struct {
	http2.WindowUpdateFrame
}

func (ev WindowUpdateFrameEvent) Type() EventType {
	return EventWindowUpdateFrame
}

func (ev WindowUpdateFrameEvent) String() string {
	return frameString(ev.Header())
}

type ContinuationFrameEvent struct {
	http2.ContinuationFrame
}

func (ev ContinuationFrameEvent) Type() EventType {
	return EventContinuationFrame
}

func (ev ContinuationFrameEvent) String() string {
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
