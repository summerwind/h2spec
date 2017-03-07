package client

import (
	"bytes"
	"encoding/binary"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func HeaderCompressionAndDecompression() *spec.ClientTestGroup {
	tg := NewTestGroup("4.3", "Header Compression and Decompression")

	// A decoding error in a header block MUST be treated as
	// a connection error (Section 5.4.1) of type COMPRESSION_ERROR.
	tg.AddTestCase(&spec.ClientTestCase{
		Desc:        "Sends invalid header block fragment",
		Requirement: "The endpoint MUST terminate the connection with a connection error of type COMPRESSION_ERROR.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			err := conn.Handshake()
			if err != nil {
				return err
			}

			req, err := conn.ReadRequest()
			if err != nil {
				return err
			}

			// Literal Header Field with Incremental Indexing without
			// Length and String segment.
			data := new(bytes.Buffer)
			data.Write([]byte("\x00\x00\x01\x01\x05"))
			binary.Write(data, binary.LittleEndian, req.StreamID)
			data.Write([]byte{0x40})

			err = conn.Send(data.Bytes())
			if err != nil {
				return err
			}

			return spec.VerifyConnectionError(conn, http2.ErrCodeCompression)
		},
	})

	return tg
}
