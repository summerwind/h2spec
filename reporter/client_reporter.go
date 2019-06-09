package reporter

import (
	"fmt"

	"github.com/summerwind/h2spec/log"
	"github.com/summerwind/h2spec/spec"
)

// SummaryForClient outputs the summary of test result that includes
// the number of passsed, skipped and failed.
func SummaryForClient(group *spec.ClientTestGroup) string {
	passed := group.PassedCount
	failed := group.FailedCount
	skipped := group.SkippedCount

	total := passed + failed + skipped
	tmp := "%d tests, %d passed, %d skipped, %d failed"
	return fmt.Sprintf(tmp, total, passed, skipped, failed)
}

func PrintSummaryForClient(group *spec.ClientTestGroup) {
	log.Println(SummaryForClient(group))
}

// PrintFailedClientTests outputs the report of failed tests.
func PrintFailedClientTests(group *spec.ClientTestGroup) {
	log.Print("Failures: \n\n")

	printClientFailed(group)
}

func printClientFailed(tg *spec.ClientTestGroup) {
	if tg.FailedCount == 0 {
		return
	}

	level := tg.Level()

	log.SetIndentLevel(level)
	log.Println(tg.Title())
	log.SetIndentLevel(level + 1)

	failed := false

	for _, tc := range tg.Tests {
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
		printClientFailed(g)
	}
}
