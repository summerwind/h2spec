package hpack

import "github.com/summerwind/h2spec/spec"

func PrimitiveTypeRepresentations() *spec.TestGroup {
	tg := NewTestGroup("5", "Primitive Type Representations")

	tg.AddTestGroup(StringLiteralRepresentation())

	return tg
}
