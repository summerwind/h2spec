package client

import "github.com/summerwind/h2spec/spec"

func ErrorHandling() *spec.ClientTestGroup {
	tg := NewTestGroup("5.4", "Error Handling")

	tg.AddTestGroup(ConnectionErrorHandling())

	return tg
}
