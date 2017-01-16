package hpack

import "github.com/summerwind/h2spec/spec"

func CompressionProcessOverview() *spec.TestGroup {
	tg := NewTestGroup("2", "Compression Process Overview")

	tg.AddTestGroup(IndexingTables())

	return tg
}
