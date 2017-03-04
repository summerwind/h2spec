package client

import "github.com/summerwind/h2spec/spec"

var key = "client"

func NewTestGroup(section string, name string) *spec.ClientTestGroup {
	return &spec.ClientTestGroup{
		Key:     key,
		Section: section,
		Name:    name,
	}
}

func Spec() *spec.ClientTestGroup {
	tg := &spec.ClientTestGroup{
		Key:  key,
		Name: "Generic tests for HTTP/2 client",
	}

	tg.AddTestGroup(StartingHTTP2())
	tg.AddTestGroup(HTTPFrames())
	tg.AddTestGroup(StreamsAndMultiplexing())
	tg.AddTestGroup(FrameDefinitions())

	return tg
}
