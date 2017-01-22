package hpack

import "github.com/summerwind/h2spec/spec"

func BinaryFormat() *spec.TestGroup {
	tg := NewTestGroup("6", "Binary Format")

	tg.AddTestGroup(IndexedHeaderFieldRepresentation())
	tg.AddTestGroup(DynamicTableSizeUpdate())

	return tg
}
