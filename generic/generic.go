package generic

import "github.com/summerwind/h2spec/spec"

var key = "generic"

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
		Name: "Generic tests for HTTP/2 server",
	}

	tg.AddTestGroup(StartingHTTP2())
	tg.AddTestGroup(StreamsAndMultiplexing())
	tg.AddTestGroup(FrameDefinitions())
	tg.AddTestGroup(HTTPMessageExchanges())
	tg.AddTestGroup(HPACK())

	return tg
}
