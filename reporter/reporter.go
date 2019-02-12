package reporter

import (
	"fmt"

	"github.com/summerwind/h2spec/log"
	"github.com/summerwind/h2spec/spec"
)

// Summary outputs the summary of test result that includes
// the number of passsed, skipped and failed.
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

// FailedTests outputs the report of failed tests.
func FailedTests(groups []*spec.TestGroup) {
	log.Print("Failures: \n\n")

	for _, tg := range groups {
		printFailed(tg)
	}
}

func printFailed(tg *spec.TestGroup) {
	if tg.FailedCount == 0 {
		return
	}

	level := tg.Level()

	log.SetIndentLevel(level)
	log.Println(tg.Title())
	log.SetIndentLevel(level + 1)

	tests := append(tg.Tests, tg.StrictTests...)
	failed := false

	for _, tc := range tests {
		if tc.Result == nil {
			continue
		}

		if tc.Result.Failed {
			tc.Result.Print()
			failed = true
		}
	}

	if failed {
		log.PrintBlankLine()
	}

	for _, g := range tg.Groups {
		printFailed(g)
	}
}
