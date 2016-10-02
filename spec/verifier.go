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

func VerifyConnectionClose(conn *Conn) error {
	var actual Event

	passed := false
	for !conn.Closed {
		actual = conn.WaitEvent()
		_, passed = actual.(EventConnectionClosed)
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

func VerifyConnectionError(conn *Conn, codes ...http2.ErrCode) error {
	var actual Event

	passed := false
	for !conn.Closed {
		actual = conn.WaitEvent()

		switch ev := actual.(type) {
		case EventConnectionClosed:
			passed = true
		case EventGoAwayFrame:
			passed = VerifyErrorCode(codes, ev.ErrCode)
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

func VerifyStreamError(conn *Conn, codes ...http2.ErrCode) error {
	var actual Event

	passed := false
	for !conn.Closed {
		actual = conn.WaitEvent()

		switch ev := actual.(type) {
		case EventConnectionClosed:
			passed = true
		case EventGoAwayFrame:
			passed = VerifyErrorCode(codes, ev.ErrCode)
		case EventRSTStreamFrame:
			passed = VerifyErrorCode(codes, ev.ErrCode)
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

func VerifyStreamClose(conn *Conn) error {
	var actual Event

	passed := false
	for !conn.Closed {
		actual = conn.WaitEvent()

		switch ev := actual.(type) {
		case EventDataFrame:
			if ev.StreamEnded() {
				passed = true
			}
		case EventHeadersFrame:
			if ev.StreamEnded() {
				passed = true
			}
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

func VerifyErrorCode(codes []http2.ErrCode, code http2.ErrCode) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}
