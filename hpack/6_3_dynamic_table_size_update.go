package hpack

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"golang.org/x/net/http2"
)

func DynamicTableSizeUpdate() *spec.TestGroup {
	tg := NewTestGroup("6.3", "Dynamic Table Size Update")

	// The new maximum size MUST be lower than or equal to the limit
	// determined by the protocol using HPACK.  A value that exceeds this
	// limit MUST be treated as a decoding error.  In HTTP/2, this limit is
	// the last value of the SETTINGS_HEADER_TABLE_SIZE parameter (see
	// Section 6.5.2 of [HTTP2]) received from the decoder and acknowledged
	// by the encoder (see Section 6.5.3 of [HTTP2]).
	tg.AddTestCase(&spec.TestCase{
		Desc:        "Sends a dynamic table size update larger than the value of SETTINGS_HEADER_TABLE_SIZE",
		Requirement: "The endpoint MUST treat this as a decoding error.",
		Run: func(c *config.Config, conn *spec.Conn) error {
			var streamID uint32 = 1

			err := conn.Handshake()
			if err != nil {
				return err
			}

			maxTableSize, ok := conn.Settings[http2.SettingHeaderTableSize]
			if !ok {
				maxTableSize = uint32(4096)
			}

			tableSize := uint64(maxTableSize + 1)
			rep := []byte{}

			// Encode to dynamic table size update.
			// This code is from golang.org/x/net/http2.
			k := uint64((1 << 5) - 1)
			if tableSize < k {
				rep = append(rep, byte(tableSize))
			} else {
				rep = append(rep, byte(k))
				tableSize -= k
				for ; tableSize >= 128; tableSize >>= 7 {
					rep = append(rep, byte(0x80|(tableSize&0x7f)))
				}
				rep = append(rep, byte(tableSize))
			}
			rep[0] |= 0x20

			headers := spec.CommonHeaders(c)
			blockFragment := conn.EncodeHeaders(headers)
			blockFragment = append(rep, blockFragment...)

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
