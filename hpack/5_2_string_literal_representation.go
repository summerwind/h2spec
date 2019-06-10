package hpack

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func StringLiteralRepresentation() *spec.TestGroup {
	tg := NewTestGroup("5.2", "String Literal Representation")

	// As the Huffman-encoded data doesn't always end at an octet boundary,
	// some padding is inserted after it, up to the next octet boundary.  To
	// prevent this padding from being misinterpreted as part of the string
	// literal, the most significant bits of the code corresponding to the
	// EOS (end-of-string) symbol are used.
	//
	// Upon decoding, an incomplete code at the end of the encoded data is
	// to be considered as padding and discarded.  A padding strictly longer
	// than 7 bits MUST be treated as a decoding error.  A padding not
	// corresponding to the most significant bits of the code for the EOS
	// symbol MUST be treated as a decoding error.  A Huffman-encoded string
	// literal containing the EOS symbol MUST be treated as a decoding
	// error.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a Huffman-encoded string literal representation with padding longer than 7 bits",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal Header Field without Indexing - New Name (x-test: test)
			// This contains an extra 1 padding octet at the end
			rep := []byte("\x00\x85\xf2\xb2\x4a\x84\xff\x84\x49\x50\x9f\xff")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(blockFragment, rep...)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: blockFragment,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeCompression)
		},
	})

	// As the Huffman-encoded data doesn't always end at an octet boundary,
	// some padding is inserted after it, up to the next octet boundary.  To
	// prevent this padding from being misinterpreted as part of the string
	// literal, the most significant bits of the code corresponding to the
	// EOS (end-of-string) symbol are used.
	//
	// Upon decoding, an incomplete code at the end of the encoded data is
	// to be considered as padding and discarded.  A padding strictly longer
	// than 7 bits MUST be treated as a decoding error.  A padding not
	// corresponding to the most significant bits of the code for the EOS
	// symbol MUST be treated as a decoding error.  A Huffman-encoded string
	// literal containing the EOS symbol MUST be treated as a decoding
	// error.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a Huffman-encoded string literal representation padded by zero",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal Header Field without Indexing - New Name (x-test: test)
			// This is padded by zero
			rep := []byte("\x00\x85\xf2\xb2\x4a\x84\xff\x83\x49\x50\x90")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(blockFragment, rep...)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: blockFragment,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeCompression)
		},
	})

	// As the Huffman-encoded data doesn't always end at an octet boundary,
	// some padding is inserted after it, up to the next octet boundary.  To
	// prevent this padding from being misinterpreted as part of the string
	// literal, the most significant bits of the code corresponding to the
	// EOS (end-of-string) symbol are used.
	//
	// Upon decoding, an incomplete code at the end of the encoded data is
	// to be considered as padding and discarded.  A padding strictly longer
	// than 7 bits MUST be treated as a decoding error.  A padding not
	// corresponding to the most significant bits of the code for the EOS
	// symbol MUST be treated as a decoding error.  A Huffman-encoded string
	// literal containing the EOS symbol MUST be treated as a decoding
	// error.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a Huffman-encoded string literal representation containing the EOS symbol",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal Header Field without Indexing - New Name (x-test: tes{EOS}t)
			// This contains a EOS symbol in the middle:
			//
			// Header:       0000 0000
			// Name Length:  1000 0101
			// Name String:  1111 0010 1011 0010 0100 1010 1000 0100 1111 1111
			// Value Length: 1000 0111
			// Value String: 0100 1001 0101 0001 1111 1111 1111 1111 1111 1111 1111 1010 0111 1111
			rep := []byte("\x00\x85\xf2\xb2\x4a\x84\xff\x87\x49\x51\xff\xff\xff\xfa\x7f")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(blockFragment, rep...)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: blockFragment,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyConnectionError(conn, http2.ErrCodeCompression)
		},
	})

	return tg
}
