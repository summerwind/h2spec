package client

import "github.com/summerwind/h2spec/spec"

func HTTPFrames() *spec.ClientTestGroup {
	tg := NewTestGroup("4", "HTTP Frames")

	tg.AddTestGroup(FrameFormat())

	return tg
}
