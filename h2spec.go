package h2spec

import (
	"fmt"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/http2"
	"github.com/summerwind/h2spec/log"
	"github.com/summerwind/h2spec/reporter"
	"github.com/summerwind/h2spec/spec"
)

func Run(c *config.Config) error {
	specs := []*spec.TestGroup{
		http2.Spec(),
	}

	start := time.Now()
	for _, s := range specs {
		s.Test(c)
	}
	end := time.Now()
	d := end.Sub(start)

	if c.DryRun {
		return nil
	}

	log.SetIndentLevel(0)
	log.Println(fmt.Sprintf("Finished in %.4f seconds", d.Seconds()))

	reporter.Summary(specs)
	reporter.FailedReport(specs)

	if c.JUnitReport != "" {
		err := reporter.JUnitReport(specs, c.JUnitReport)
		if err != nil {
			return err
		}
	}

	return nil
}
