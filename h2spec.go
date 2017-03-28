package h2spec

import (
	"fmt"
	"time"

	"github.com/summerwind/h2spec/client"
	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/generic"
	"github.com/summerwind/h2spec/hpack"
	"github.com/summerwind/h2spec/http2"
	"github.com/summerwind/h2spec/log"
	"github.com/summerwind/h2spec/reporter"
	"github.com/summerwind/h2spec/spec"
)

func Run(c *config.Config) (bool, error) {
	total := 0
	success := true

	specs := []*spec.TestGroup{
		generic.Spec(),
		http2.Spec(),
		hpack.Spec(),
	}

	start := time.Now()
	for _, s := range specs {
		s.Test(c)

		if s.FailedCount > 0 {
			success = false
		}

		total += s.FailedCount
		total += s.SkippedCount
		total += s.PassedCount
	}
	end := time.Now()
	d := end.Sub(start)

	if c.DryRun {
		return true, nil
	}

	if total == 0 {
		log.SetIndentLevel(0)
		log.Println("No matched tests found.")
		return true, nil
	}

	if !success {
		log.SetIndentLevel(0)
		reporter.FailedTests(specs)
	}

	log.SetIndentLevel(0)
	log.Println(fmt.Sprintf("Finished in %.4f seconds", d.Seconds()))
	reporter.Summary(specs)

	if c.JUnitReport != "" {
		err := reporter.JUnitReport(specs, c.JUnitReport)
		if err != nil {
			return false, err
		}
	}

	return success, nil
}

func RunClientSpec(c *config.Config) error {
	s := client.Spec()

	server, err := spec.Listen(c, s)
	if err != nil {
		return err
	}

	if !c.IsBrowserMode() {
		start := time.Now()
		s.Test(c)
		end := time.Now()
		d := end.Sub(start)

		if s.FailedCount > 0 {
			log.SetIndentLevel(0)
			reporter.PrintFailedClientTests(s)
		}

		log.SetIndentLevel(0)
		log.Println(fmt.Sprintf("Finished in %.4f seconds", d.Seconds()))
		reporter.PrintSummaryForClient(s)
	} else {
		// Block running
		log.Println("--exec is not defined, enable BROWSER mode")

		reportServer := reporter.NewWebReportServer(c, s)
		log.Println(reportServer.RunForever())
	}

	defer server.Close()

	return nil
}
