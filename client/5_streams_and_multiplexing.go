package client

import "github.com/summerwind/h2spec/spec"

func StreamsAndMultiplexing() *spec.ClientTestGroup {
	tg := NewTestGroup("5", "Streams and Multiplexing")

	tg.AddTestGroup(StreamStates())
	tg.AddTestGroup(ErrorHandling())
	tg.AddTestGroup(ExtendingHTTP2())

	return tg
}
