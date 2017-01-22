package http2

import "github.com/summerwind/h2spec/spec"

func FrameDefinitions() *spec.TestGroup {
	tg := NewTestGroup("6", "Frame Definitions")

	tg.AddTestGroup(Data())
	tg.AddTestGroup(Headers())
	tg.AddTestGroup(Priority())
	tg.AddTestGroup(RSTStream())
	tg.AddTestGroup(Settings())
	tg.AddTestGroup(Ping())
	tg.AddTestGroup(GoAway())
	tg.AddTestGroup(WindowUpdate())
	tg.AddTestGroup(Continuation())

	return tg
}
