package http2

import "github.com/summerwind/h2spec/spec"

func StartingHTTP2() *spec.TestGroup {
	tg := NewTestGroup("3", "Starting HTTP/2")

	tg.AddTestGroup(HTTP2ConnectionPreface())

	return tg
}
