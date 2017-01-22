package http2

import "github.com/summerwind/h2spec/spec"

var key = "http2"

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
		Name: "Hypertext Transfer Protocol Version 2 (HTTP/2)",
	}

	tg.AddTestGroup(StartingHTTP2())
	tg.AddTestGroup(HTTPFrames())
	tg.AddTestGroup(StreamsAndMultiplexing())
	tg.AddTestGroup(FrameDefinitions())
	tg.AddTestGroup(ErrorCodes())
	tg.AddTestGroup(HTTPMessageExchanges())

	return tg
}
