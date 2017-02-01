package generic

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func HPACK() *spec.TestGroup {
	tg := NewTestGroup("5", "HPACK")

	// RFC 7541, 6.1:
	// An indexed header field representation identifies an entry in either
	// the static table or the dynamic table (see Section 2.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a indexed header field representation",
		Requirement: "The endpoint MUST accept indexed header field representation",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Indexed header field representation
			// (user-agent: )
			rep := []byte("\xba")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.1:
	// A literal header field with incremental indexing representation
	// results in appending a header field to the decoded header list and
	// inserting it as a new entry into the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field with incremental indexing - indexed name",
		Requirement: "The endpoint MUST accept literal header field with incremental indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field with incremental indexing - indexed name
			// (user-agent: h2spec)
			rep := []byte("\x40\x87\xb5\x05\xb1\x61\xcc\x5a\x93\x84\x9c\x48\xac\xa4")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.1:
	// A literal header field with incremental indexing representation
	// results in appending a header field to the decoded header list and
	// inserting it as a new entry into the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field with incremental indexing - indexed name (with Huffman coding)",
		Requirement: "The endpoint MUST accept literal header field with incremental indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field with incremental indexing - indexed name
			// (user-agent: h2spec)
			rep := []byte("\x40\x0a\x75\x73\x65\x72\x2d\x61\x67\x65\x6e\x74\x06\x68\x32\x73\x70\x65\x63")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.1:
	// A literal header field with incremental indexing representation
	// results in appending a header field to the decoded header list and
	// inserting it as a new entry into the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field with incremental indexing - new name",
		Requirement: "The endpoint MUST accept literal header field with incremental indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field with incremental indexing - new name
			// (x-test: h2spec)
			rep := []byte("\x40\x06\x78\x2d\x74\x65\x73\x74\x06\x68\x32\x73\x70\x65\x63")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.1:
	// A literal header field with incremental indexing representation
	// results in appending a header field to the decoded header list and
	// inserting it as a new entry into the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field with incremental indexing - new name (with Huffman coding)",
		Requirement: "The endpoint MUST accept literal header field with incremental indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field with incremental indexing - new name
			// (x-test: h2spec)
			rep := []byte("\x40\x85\xf2\xb2\x4a\x84\xff\x84\x9c\x48\xac\xa4")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.2:
	// A literal header field without indexing representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field without indexing - indexed name",
		Requirement: "The endpoint MUST accept literal header field without indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field without indexing - indexed name
			// (user-agent: h2spec)
			rep := []byte("\x00\x0a\x75\x73\x65\x72\x2d\x61\x67\x65\x6e\x74\x06\x68\x32\x73\x70\x65\x63")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.2:
	// A literal header field without indexing representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field without indexing - indexed name (with Huffman coding)",
		Requirement: "The endpoint MUST accept literal header field without indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field without indexing - indexed name
			// (user-agent: h2spec)
			rep := []byte("\x00\x87\xb5\x05\xb1\x61\xcc\x5a\x93\x84\x9c\x48\xac\xa4")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.2:
	// A literal header field without indexing representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field without indexing - new name",
		Requirement: "The endpoint MUST accept literal header field without indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field without indexing - new name
			// (x-test: h2spec)
			rep := []byte("\x00\x06\x78\x2d\x74\x65\x73\x74\x06\x68\x32\x73\x70\x65\x63")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.2:
	// A literal header field without indexing representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field without indexing - new name (huffman encoded)",
		Requirement: "The endpoint MUST accept literal header field without indexing",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field without indexing - new name
			// (x-test: h2spec)
			rep := []byte("\x00\x85\xf2\xb2\x4a\x84\xff\x84\x9c\x48\xac\xa4")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.3:
	// A literal header field never-indexed representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.  Intermediaries MUST use the same representation
	// for encoding this header field.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field never indexed - indexed name",
		Requirement: "The endpoint MUST accept literal header field never indexed",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field never indexed - indexed name
			// (user-agent: h2spec)
			rep := []byte("\x10\x0a\x75\x73\x65\x72\x2d\x61\x67\x65\x6e\x74\x06\x68\x32\x73\x70\x65\x63")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.3:
	// A literal header field never-indexed representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.  Intermediaries MUST use the same representation
	// for encoding this header field.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field never indexed - indexed name (huffman encoded)",
		Requirement: "The endpoint MUST accept literal header field never indexed",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field never indexed - indexed name
			// (user-agent: h2spec)
			rep := []byte("\x10\x87\xb5\x05\xb1\x61\xcc\x5a\x93\x84\x9c\x48\xac\xa4")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.3:
	// A literal header field never-indexed representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.  Intermediaries MUST use the same representation
	// for encoding this header field.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field never indexed - new name",
		Requirement: "The endpoint MUST accept literal header field never indexed",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field never indexed - new name
			// (x-test: h2spec)
			rep := []byte("\x10\x85\xf2\xb2\x4a\x84\xff\x84\x9c\x48\xac\xa4")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.2.3:
	// A literal header field never-indexed representation results in
	// appending a header field to the decoded header list without altering
	// the dynamic table.  Intermediaries MUST use the same representation
	// for encoding this header field.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a literal header field never indexed - new name (huffman encoded)",
		Requirement: "The endpoint MUST accept literal header field never indexed",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Literal header field never indexed - new name
			// (x-test: h2spec)
			rep := []byte("\x10\x85\xf2\xb2\x4a\x84\xff\x84\x9c\x48\xac\xa4")

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

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 6.3:
	// A dynamic table size update signals a change to the size of the
	// dynamic table.
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a dynamic table size update",
		Requirement: "The endpoint MUST accept dynamic table size update",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// Change encoder's table size to 128. This sends dynamic table
			// size update with value 128.
			conn.SetMaxDynamicTableSize(128)

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: blockFragment,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	// RFC 7541, 4.2:
	// A change in the maximum size of the dynamic table is signaled via a
	// dynamic table size update (see Section 6.3).  This dynamic table size
	// update MUST occur at the beginning of the first header block
	// following the change to the dynamic table size.  In HTTP/2, this
	// follows a settings acknowledgment (see Section 6.5.3 of [HTTP2]).
	//
	// Multiple updates to the maximum table size can occur between the
	// transmission of two header blocks.  In the case that this size is
	// changed more than once in this interval, the smallest maximum table
	// size that occurs in that interval MUST be signaled in a dynamic table
	// size update.  The final maximum size is always signaled, resulting in
	// at most two dynamic table size updates.  This ensures that the
	// decoder is able to perform eviction based on reductions in dynamic
	// table size (see Section 4.3).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends multiple dynamic table size update",
		Requirement: "The endpoint MUST accept multiple dynamic table size update",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			// 2 Dynamic Table Size Updates, 128 and 4096.
			tableSizeUpdate := []byte("\x3f\x61\x3f\xe1\x1f")

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(tableSizeUpdate, blockFragment...)

			hp := http2.HeadersFrameParam{
				StreamID:      streamID,
				EndStream:     true,
				EndHeaders:    true,
				BlockFragment: blockFragment,
			}
			conn.WriteHeaders(hp)

			return spec.VerifyHeadersFrame(conn, streamID)
		},
	})

	return tg
}
