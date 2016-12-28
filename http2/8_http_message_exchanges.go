package http2

import "github.com/summerwind/h2spec/spec"

func HTTPMessageExchanges() *spec.TestGroup {
	tg := NewTestGroup("8", "HTTP Message Exchanges")

	tg.AddTestGroup(HTTPRequestResponseExchange())

	return tg
}
