package spec

import (
	"fmt"

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

// VerifyFrameType verifies whether a frame with specified type
// has received.
func VerifyFrameType(conn *Conn, frameTypes ...http2.FrameType) error {
	var actual Event

	passed := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case TimeoutEvent:
			if actual == nil {
				actual = event
			}
		default:
			actual = ev

			ef, ok := event.(EventFrame)
			if ok {
				for _, ft := range frameTypes {
					if ef.Header().Type == ft {
						passed = true
					}
				}
			}
		}

		if passed {
			break
		}
	}

	if !passed {
		expected := []string{}
		for _, ft := range frameTypes {
			expected = append(expected, fmt.Sprintf("%s frame", ft))
		}

		return &TestError{
			Expected: expected,
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
