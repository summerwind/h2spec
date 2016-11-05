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
		event := conn.WaitEvent()

		switch ev := event.(type) {
		case EventConnectionClosed:
			passed = true
		case EventTimeout:
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

func VerifyConnectionError(conn *Conn, codes ...http2.ErrCode) error {
	var actual Event

	passed := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case EventConnectionClosed:
			passed = true
		case EventGoAwayFrame:
			passed = VerifyErrorCode(codes, event.ErrCode)
		case EventTimeout:
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

func VerifyStreamError(conn *Conn, codes ...http2.ErrCode) error {
	var actual Event

	passed := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case EventConnectionClosed:
			passed = true
		case EventGoAwayFrame:
			passed = VerifyErrorCode(codes, event.ErrCode)
		case EventRSTStreamFrame:
			passed = VerifyErrorCode(codes, event.ErrCode)
		case EventTimeout:
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

func VerifyStreamClose(conn *Conn) error {
	var actual Event

	passed := false
	for !conn.Closed {
		ev := conn.WaitEvent()

		switch event := ev.(type) {
		case EventDataFrame:
			if event.StreamEnded() {
				passed = true
			}
		case EventHeadersFrame:
			if event.StreamEnded() {
				passed = true
			}
		case EventTimeout:
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

func VerifyErrorCode(codes []http2.ErrCode, code http2.ErrCode) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}
