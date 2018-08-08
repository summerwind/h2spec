package spec

import (
	"fmt"
	"reflect"

	"golang.org/x/net/http2"
)

const (
	ExpectedConnectionClosed = "Connection closed"
	ExpectedStreamClosed     = "Stream closed"
	ExpectedGoAwayFrame      = "GOAWAY Frame (Error Code: %s)"
	ExpectedRSTStreamFrame   = "RST_STREAM Frame (Error Code: %s)"
)

// VerifyConnectionClose verifies whether the connection was closed.
func VerifyConnectionClose(conn *Conn) error {
	var actual Event

	passed := false
	for !conn.Closed {
		event := conn.WaitEvent()

		switch ev := event.(type) {
		case ConnectionClosedEvent:
			passed = true
		case TimeoutEvent:
			if actual == nil {
				actual = ev
			}
		default:
			actual = ev
		}

		if passed {
			break
		}
	}

	if !passed {
		return &TestError{
			Expected: []string{ExpectedConnectionClosed},
			Actual:   actual.String(),
		}
	}

	return nil
}

// VerifyConnectionError verifies whether a connection error of HTTP/2
// has occurred.
func VerifyConnectionError(conn *Conn, codes ...http2.ErrCode) error {
	var actual Event

	passed := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case ConnectionClosedEvent:
			passed = true
		case GoAwayFrameEvent:
			passed = VerifyErrorCode(codes, event.ErrCode)
		case TimeoutEvent:
			if actual == nil {
				actual = event
			}
		default:
			actual = event
		}

		if passed {
			break
		}
	}

	if !passed {
		expected := []string{}
		for _, code := range codes {
			expected = append(expected, fmt.Sprintf(ExpectedGoAwayFrame, code))
		}
		expected = append(expected, ExpectedConnectionClosed)

		return &TestError{
			Expected: expected,
			Actual:   actual.String(),
		}
	}

	return nil
}

// VerifyStreamError verifies whether a stream error of HTTP/2
// has occurred.
func VerifyStreamError(conn *Conn, codes ...http2.ErrCode) error {
	var actual Event

	passed := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case ConnectionClosedEvent:
			passed = true
		case GoAwayFrameEvent:
			passed = VerifyErrorCode(codes, event.ErrCode)
		case RSTStreamFrameEvent:
			passed = VerifyErrorCode(codes, event.ErrCode)
		case TimeoutEvent:
			if actual == nil {
				actual = event
			}
		default:
			actual = event
		}

		if passed {
			break
		}
	}

	if !passed {
		expected := []string{}
		for _, code := range codes {
			expected = append(expected, fmt.Sprintf(ExpectedGoAwayFrame, code))
			expected = append(expected, fmt.Sprintf(ExpectedRSTStreamFrame, code))
		}
		expected = append(expected, ExpectedConnectionClosed)

		return &TestError{
			Expected: expected,
			Actual:   actual.String(),
		}
	}

	return nil
}

// VerifyStreamClose verifies whether a stream close of HTTP/2
// has occurred.
func VerifyStreamClose(conn *Conn) error {
	var actual Event

	passed := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case DataFrameEvent:
			if event.StreamEnded() {
				passed = true
			}
		case HeadersFrameEvent:
			if event.StreamEnded() {
				passed = true
			}
		case RSTStreamFrameEvent:
			if event.ErrCode == http2.ErrCodeNo {
				passed = true
			}
		case TimeoutEvent:
			if actual == nil {
				actual = event
			}
		default:
			actual = event
		}

		if passed {
			break
		}
	}

	if !passed {
		return &TestError{
			Expected: []string{ExpectedStreamClosed},
			Actual:   actual.String(),
		}
	}

	return nil
}

// VerifyHeadersFrame verifies whether a HEADERS frame with specified
// stream ID has received.
func VerifyHeadersFrame(conn *Conn, streamID uint32) error {
	actual, passed := conn.WaitEventByType(EventHeadersFrame)
	switch event := actual.(type) {
	case HeadersFrameEvent:
		passed = (event.Header().StreamID == streamID)
	default:
		passed = false
	}

	if !passed {
		expected := []string{
			fmt.Sprintf("HEADERS Frame (stream_id:%d)", streamID),
		}

		return &TestError{
			Expected: expected,
			Actual:   actual.String(),
		}
	}

	return nil
}

// VerifySettingsFrameWithAck verifies whether a SETTINGS frame with
// ACK flag has received.
func VerifySettingsFrameWithAck(conn *Conn) error {
	actual, passed := conn.WaitEventByType(EventSettingsFrame)
	switch event := actual.(type) {
	case SettingsFrameEvent:
		passed = event.IsAck()
	default:
		passed = false
	}

	if !passed {
		expected := []string{
			"SETTINGS Frame (length:0, flags:0x01, stream_id:0)",
		}

		return &TestError{
			Expected: expected,
			Actual:   actual.String(),
		}
	}

	return nil
}

// VerifyPingFrameWithAck verifies whether a PING frame with ACK flag
// has received.
func VerifyPingFrameWithAck(conn *Conn, data [8]byte) error {
	actual, passed := conn.WaitEventByType(EventPingFrame)
	switch event := actual.(type) {
	case PingFrameEvent:
		passed = event.IsAck() && reflect.DeepEqual(event.Data, data)
	default:
		passed = false
	}

	if !passed {
		var actualStr string

		expected := []string{
			fmt.Sprintf("PING Frame (length:8, flags:0x01, stream_id:0, opaque_data:%s)", data),
		}

		f, ok := actual.(PingFrameEvent)
		if ok {
			header := f.Header()
			actualStr = fmt.Sprintf(
				"PING Frame ((length:%d, flags:0x%02x, stream_id:%d, opaque_data: %s)",
				header.Length,
				header.Flags,
				header.StreamID,
				f.Data,
			)
		} else {
			actualStr = actual.String()
		}

		return &TestError{
			Expected: expected,
			Actual:   actualStr,
		}
	}

	return nil
}

// VerifyPingFrameOrConnectionClose verifies whether a PING frame with
// ACK flag has received or the connection was closed.
func VerifyPingFrameOrConnectionClose(conn *Conn, data [8]byte) error {
	var actual Event

	passed := false
	for !conn.Closed {
		event := conn.WaitEvent()

		switch ev := event.(type) {
		case ConnectionClosedEvent:
			passed = true
		case PingFrameEvent:
			passed = ev.IsAck() && reflect.DeepEqual(ev.Data, data)
		case TimeoutEvent:
			if actual == nil {
				actual = ev
			}
		default:
			actual = ev
		}

		if passed {
			break
		}
	}

	if !passed {
		var actualStr string

		expected := []string{
			ExpectedConnectionClosed,
			fmt.Sprintf("PING Frame (length:8, flags:0x01, stream_id:0, opaque_data:%s)", data),
		}

		f, ok := actual.(PingFrameEvent)
		if ok {
			header := f.Header()
			actualStr = fmt.Sprintf(
				"PING Frame ((length:%d, flags:0x%02x, stream_id:%d, opaque_data: %s)",
				header.Length,
				header.Flags,
				header.StreamID,
				f.Data,
			)
		} else {
			actualStr = actual.String()
		}

		return &TestError{
			Expected: expected,
			Actual:   actualStr,
		}
	}

	return nil
}

// VerifyEventType verifies whether a frame with specified type
// has received.
func VerifyEventType(conn *Conn, et EventType) error {
	var actual Event

	actual, passed := conn.WaitEventByType(et)

	if !passed {
		return &TestError{
			Expected: []string{et.String()},
			Actual:   actual.String(),
		}
	}

	return nil
}

// VerifyErrorCode verifies whether the specified error code is
// the expected error code.
func VerifyErrorCode(codes []http2.ErrCode, code http2.ErrCode) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}
