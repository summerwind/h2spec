package client

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StreamStates() *spec.ClientTestGroup {
	tg := NewTestGroup("5.1", "Stream States")

	// idle:
	// Receiving any frame other than HEADERS or PRIORITY on a stream
	// in this state MUST be treated as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "idle: Sends a DATA frame",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteData(2, true, []byte("test"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// idle:
	// Receiving any frame other than HEADERS or PRIORITY on a stream
	// in this state MUST be treated as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "idle: Sends a RST_STREAM frame",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteRSTStream(2, http2.ErrCodeCancel)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// idle:
	// Receiving any frame other than HEADERS or PRIORITY on a stream
	// in this state MUST be treated as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "idle: Sends a WINDOW_UPDATE frame",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			conn.WriteWindowUpdate(2, 100)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// idle:
	// Receiving any frame other than HEADERS or PRIORITY on a stream
	// in this state MUST be treated as a connection error
	// (Section 5.4.1) of type PROTOCOL_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "idle: Sends a CONTINUATION frame",
		Requirement: "The endpoint MUST treat this as a connection error of type PROTOCOL_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			conn.WriteContinuation(2, true, blockFragment)

			return spec.VerifyConnectionError(conn, http2.ErrCodeProtocol)
		},
	})

	// closed:
	// An endpoint that receives any frame other than PRIORITY after
	// receiving a RST_STREAM MUST treat that as a stream error
	// (Section 5.4.2) of type STREAM_CLOSED.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "closed: Sends a DATA frame after sending RST_STREAM frame",
		Requirement: "The endpoint MUST treat this as a stream error of type STREAM_CLOSED.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			conn.WriteRSTStream(req.StreamID, http2.ErrCodeCancel)

			conn.WriteData(req.StreamID, true, []byte("test"))

			return spec.VerifyStreamError(conn, http2.ErrCodeStreamClosed)
		},
	})

	// closed:
	// An endpoint that receives any frame other than PRIORITY after
	// receiving a RST_STREAM MUST treat that as a stream error
	// (Section 5.4.2) of type STREAM_CLOSED.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "closed: Sends a HEADERS frame after sending RST_STREAM frame",
		Requirement: "The endpoint MUST treat this as a stream error of type STREAM_CLOSED.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp1 := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp1)

			conn.WriteRSTStream(req.StreamID, http2.ErrCodeCancel)

			hp2 := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp2)

			return spec.VerifyStreamError(conn, http2.ErrCodeStreamClosed)
		},
	})

	// closed:
	// An endpoint that receives any frame other than PRIORITY after
	// receiving a RST_STREAM MUST treat that as a stream error
	// (Section 5.4.2) of type STREAM_CLOSED.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "closed: Sends a CONTINUATION frame after sending RST_STREAM frame",
		Requirement: "The endpoint MUST treat this as a stream error of type STREAM_CLOSED.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     false,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			conn.WriteRSTStream(req.StreamID, http2.ErrCodeCancel)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, true, conn.EncodeHeaders(dummyHeaders))

			codes := []http2.ErrCode{
				http2.ErrCodeStreamClosed,
				http2.ErrCodeProtocol,
			}
			return spec.VerifyStreamError(conn, codes...)
		},
	})

	// closed:
	// An endpoint that receives any frames after receiving a frame
	// with the END_STREAM flag set MUST treat that as a connection
	// error (Section 6.4.1) of type STREAM_CLOSED.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "closed: Sends a DATA frame",
		Requirement: "The endpoint MUST treat this as a connection error of type STREAM_CLOSED.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			conn.WriteData(req.StreamID, true, []byte("test"))

			return spec.VerifyConnectionError(conn, http2.ErrCodeStreamClosed)
		},
	})

	// closed:
	// An endpoint that receives any frames after receiving a frame
	// with the END_STREAM flag set MUST treat that as a connection
	// error (Section 6.4.1) of type STREAM_CLOSED.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "closed: Sends a HEADERS frame",
		Requirement: "The endpoint MUST treat this as a connection error of type STREAM_CLOSED.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeStreamClosed)
		},
	})

	// closed:
	// An endpoint that receives any frames after receiving a frame
	// with the END_STREAM flag set MUST treat that as a connection
	// error (Section 6.4.1) of type STREAM_CLOSED.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "closed: Sends a CONTINUATION frame",
		Requirement: "The endpoint MUST treat this as a connection error of type STREAM_CLOSED.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			headers := spec.CommonRespHeaders(c)
			hp := http2.HeadersFrameParam{
				StreamID:      req.StreamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: conn.EncodeHeaders(headers),
			}
			conn.WriteHeaders(hp)

			dummyHeaders := spec.DummyRespHeaders(c, 1)
			conn.WriteContinuation(req.StreamID, true, conn.EncodeHeaders(dummyHeaders))

			codes := []http2.ErrCode{
				http2.ErrCodeStreamClosed,
				http2.ErrCodeProtocol,
			}
			return spec.VerifyConnectionError(conn, codes...)
		},
	})

	tg.AddTestGroup(StreamIdentifiers())

	return tg
}
