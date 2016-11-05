package reporter

import (
	"fmt"

	"github.com/summerwind/h2spec/log"
	"github.com/summerwind/h2spec/spec"
)

func Summary(groups []*spec.TestGroup) {
	s := aggregateSummary(groups)
	tmp := "%d tests, %d passed, %d skipped, %d failed"
	log.Println(fmt.Sprintf(tmp, s["total"], s["passed"], s["skipped"], s["failed"]))
}

func aggregateSummary(groups []*spec.TestGroup) map[string]int {
	data := map[string]int{
		"total":   0,
		"passed":  0,
		"failed":  0,
		"skipped": 0,
	}

	for _, tg := range groups {
		tests := append(tg.Tests, tg.StrictTests...)
		for _, tc := range tests {
			if tc.Result == nil {
				continue
			}

			data["total"] += 1
			if tc.Result.Failed {
				data["failed"] += 1
			} else if tc.Result.Skipped {
				data["skipped"] += 1
			} else {
				data["passed"] += 1
			}
		}

		if tg.Groups != nil {
			d := aggregateSummary(tg.Groups)
			data["total"] += d["total"]
			data["failed"] += d["failed"]
			data["skipped"] += d["skipped"]
			data["passed"] += d["passed"]
		}
	}

	return data
}

func FailedReport(groups []*spec.TestGroup) {
	log.Println("\n--------------------\n")
	log.Println("Failed tests: \n")

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
