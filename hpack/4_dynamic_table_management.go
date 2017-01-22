package hpack

import "github.com/summerwind/h2spec/spec"

func DynamicTableManagement() *spec.TestGroup {
	tg := NewTestGroup("4", "Dynamic Table Management")

	tg.AddTestGroup(MaximumTableSize())

	return tg
}
