package reporter

import (
	"fmt"

	"github.com/summerwind/h2spec/log"
	"github.com/summerwind/h2spec/spec"
)

func Summary(groups []*spec.TestGroup) {
	var passed, failed, skipped, total int

	for _, tg := range groups {
		passed += tg.PassedCount
		failed += tg.FailedCount
		skipped += tg.SkippedCount
	}

	total = passed + failed + skipped
	tmp := "%d tests, %d passed, %d skipped, %d failed"
	log.Println(fmt.Sprintf(tmp, total, passed, skipped, failed))
}

func FailedReport(groups []*spec.TestGroup) {
	failed := false

	for _, tg := range groups {
		if tg.FailedCount > 0 {
			failed = true
		}
	}

	if !failed {
		return
	}

	log.Println("Failures: \n")

	for _, tg := range groups {
		log.Println(tg.Title())
		printFailed(tg.Groups)
	}
}

func printFailed(groups []*spec.TestGroup) {
	for _, tg := range groups {
		failed := false

		tests := append(tg.Tests, tg.StrictTests...)
		for _, tc := range tests {
			if tc.Result == nil {
				continue
			}

			if tc.Result.Failed {
				log.SetIndentLevel(1)
				log.Println(tc.Parent.Title())
				log.SetIndentLevel(2)
				tc.Result.Print()
				failed = true
			}
		}

		if failed {
			log.PrintBlankLine()
		}

		if tg.Groups != nil {
			printFailed(tg.Groups)
		}
	}
}
