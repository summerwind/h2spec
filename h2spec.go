package h2spec

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec"
	"github.com/summerwind/h2spec/spec/http2"
	"github.com/summerwind/h2spec/spec/log"
)

func Run(c *config.Config) error {
	specs := []*spec.TestGroup{
		http2.Spec(),
	}

	results := []*spec.TestResult{}
	start := time.Now()
	for _, s := range specs {
		results = append(results, s.Test(c)...)
	}
	end := time.Now()
	d := end.Sub(start)

	if c.DryRun {
		return nil
	}

	log.SetIndentLevel(0)
	log.Info(fmt.Sprintf("Finished in %.4f seconds", d.Seconds()))

	s := summary(results)
	tmp := "%d tests, %d passed, %d skipped, %d failed"
	log.Info(fmt.Sprintf(tmp, s["total"], s["passed"], s["skipped"], s["failed"]))

	if c.JUnitReport != "" {
		reporter := NewJUnitReporter()
		report, err := reporter.Export(results)
		if err != nil {
			return err
		}
		ioutil.WriteFile(c.JUnitReport, []byte(report), os.ModePerm)
	}

	return nil
}

func summary(results []*spec.TestResult) map[string]int {
	data := map[string]int{
		"total":   0,
		"passed":  0,
		"failed":  0,
		"skipped": 0,
	}

	for _, result := range results {
		data["total"] += 1
		if result.Failed() {
			data["failed"] += 1
		} else if result.Skipped() {
			data["skipped"] += 1
		} else {
			data["passed"] += 1
		}
	}

	return data
}
