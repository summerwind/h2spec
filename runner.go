package h2spec

import (
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"github.com/summerwind/h2spec/spec/http2"
)

type Runner struct {
	Specs  []*spec.TestGroup
	Config *config.Config
}

func (r *Runner) Run() ([]*spec.TestResult, error) {
	for _, s := range r.Specs {
		s.Test(r.Config)
	}
	return nil, nil
}

func NewRunner(c *config.Config) *Runner {
	return &Runner{
		Specs: []*spec.TestGroup{
			http2.Spec(),
		},
		Config: c,
	}
}
