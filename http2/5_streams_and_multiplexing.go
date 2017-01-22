package http2

import "github.com/summerwind/h2spec/spec"

func StreamsAndMultiplexing() *spec.TestGroup {
	tg := NewTestGroup("5", "Streams and Multiplexing")

	tg.AddTestGroup(StreamStates())
	tg.AddTestGroup(StreamPriority())
	tg.AddTestGroup(ErrorHandling())
	tg.AddTestGroup(ExtendingHTTP2())

	return tg
}
