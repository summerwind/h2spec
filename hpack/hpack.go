package hpack

import "github.com/summerwind/h2spec/spec"

var key = "hpack"

func NewTestGroup(section string, name string) *spec.TestGroup {
	return &spec.TestGroup{
		Key:     key,
		Section: section,
		Name:    name,
	}
}

func Spec() *spec.TestGroup {
	tg := &spec.TestGroup{
		Key:  key,
		Name: "HPACK: Header Compression for HTTP/2",
	}

	tg.AddTestGroup(CompressionProcessOverview())
	tg.AddTestGroup(DynamicTableManagement())
	tg.AddTestGroup(PrimitiveTypeRepresentations())
	tg.AddTestGroup(BinaryFormat())

	return tg
}
