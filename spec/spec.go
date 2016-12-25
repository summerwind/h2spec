package spec

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
)

var (
	ErrTimeout = errors.New("Timeout")
	ErrSkipped = errors.New("Skipped")
)

type TestGroup struct {
	Key         string
	Section     string
	Name        string
	Strict      bool
	Parent      *TestGroup
	Groups      []*TestGroup
	Tests       []*TestCase
	StrictTests []*TestCase

	PassedCount  int
	FailedCount  int
	SkippedCount int
}

func (tg *TestGroup) IsRoot() bool {
	return tg.Parent == nil
}

func (tg *TestGroup) ID() string {
	if tg.IsRoot() {
		return tg.Key
	}

	return fmt.Sprintf("%s/%s", tg.Key, tg.Section)
}

func (tg *TestGroup) Title() string {
	if tg.IsRoot() {
		return fmt.Sprintf("%s", tg.Name)
	} else {
		return fmt.Sprintf("%s. %s", tg.Section, tg.Name)
	}
}

func (tg *TestGroup) Level() int {
	if tg.IsRoot() {
		return 0
	}

	return strings.Count(tg.Section, ".") + 1
}

func (tg *TestGroup) Test(c *config.Config) {
	level := tg.Level()

	if tg.Strict && !c.Strict {
		return
	}

	mode := c.RunMode(tg.ID())
	if mode == config.RunModeNone {
		return
	}

	log.SetIndentLevel(level)
	log.Println(tg.Title())
	log.SetIndentLevel(level + 1)

	tests := append(tg.Tests, tg.StrictTests...)
	tested := false

	for i, tc := range tests {
		seq := i + 1

		err := tc.Test(c, seq)
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			os.Exit(1)
		}

		if tc.Result != nil {
			if tc.Result.Failed {
				tg.FailedCount += 1
			} else if tc.Result.Skipped {
				tg.SkippedCount += 1
			} else {
				tg.PassedCount += 1
			}
		}

		tested = true
	}

	if tested {
		log.PrintBlankLine()
	}

	for _, g := range tg.Groups {
		g.Test(c)
		tg.FailedCount += g.FailedCount
		tg.SkippedCount += g.SkippedCount
		tg.PassedCount += g.PassedCount
	}
}

func (tg *TestGroup) AddTestGroup(stg *TestGroup) {
	stg.Parent = tg
	stg.Strict = tg.Strict
	tg.Groups = append(tg.Groups, stg)
}

func (tg *TestGroup) AddTestCase(tc *TestCase) {
	tc.Parent = tg
	if tg.Strict {
		tc.Strict = true
		tg.StrictTests = append(tg.StrictTests, tc)
	} else {
		tg.Tests = append(tg.Tests, tc)
	}
}

type TestCase struct {
	Desc        string
	Requirement string
	Strict      bool
	Parent      *TestGroup
	Result      *TestResult
	Run         func(c *config.Config, conn *Conn) error
}

func (tc *TestCase) Test(c *config.Config, seq int) error {
	if c.DryRun {
		msg := fmt.Sprintf("%s %s", seqStr(seq), tc.Desc)
		log.Println(msg)
		return nil
	}

	if tc.Strict && !c.Strict {
		return nil
	}

	mode := c.RunMode(fmt.Sprintf("%s/%d", tc.Parent.ID(), seq))
	if mode == config.RunModeNone {
		return nil
	}

	if !c.Verbose {
		msg := gray(fmt.Sprintf("  %s %s", seqStr(seq), tc.Desc))
		log.Print(msg)
	}

	conn, err := Dial(c)
	if err != nil {
		msg := red(fmt.Sprintf("%s %s %s", "×", seqStr(seq), tc.Desc))
		log.ResetLine()
		log.Println(msg)
		return err
	}
	defer conn.Close()

	start := time.Now()
	err = tc.Run(c, conn)
	end := time.Now()

	log.ResetLine()

	tr := NewTestResult(tc, seq, err, end.Sub(start))
	tr.Print()
	tc.Result = tr

	return nil
}

type TestError struct {
	Expected []string
	Actual   string
}

func (e TestError) Error() string {
	return fmt.Sprintf("%s\n%s", strings.Join(e.Expected, "\n"), e.Actual)
}

type TestResult struct {
	TestCase *TestCase
	Sequence int
	Error    error
	Duration time.Duration

	Skipped bool
	Failed  bool
}

func NewTestResult(tc *TestCase, seq int, err error, d time.Duration) *TestResult {
	skipped := false
	failed := false

	if err != nil {
		if err == ErrSkipped {
			skipped = true
		} else {
			failed = true
		}
	}

	tr := TestResult{
		TestCase: tc,
		Sequence: seq,
		Error:    err,
		Duration: d,
		Skipped:  skipped,
		Failed:   failed,
	}

	return &tr
}

func (tr *TestResult) Print() {
	tc := tr.TestCase
	desc := tc.Desc
	seq := seqStr(tr.Sequence)

	if tr.Skipped {
		log.Println(cyan(fmt.Sprintf("%s %s", seq, desc)))
		return
	}

	if !tr.Failed {
		log.Println(fmt.Sprintf("%s %s %s", green("✔"), gray(seq), gray(desc)))
		return
	}

	log.Println(red(fmt.Sprintf("%s %s %s", "×", seq, desc)))
	err, ok := tr.Error.(*TestError)
	if ok {
		level := log.IndentLevel
		log.SetIndentLevel(level + 1)
		defer func() {
			log.SetIndentLevel(level)
		}()

		log.Println(red(fmt.Sprintf("-> %s", tc.Requirement)))
		label := "Expected: "
		for i, ex := range err.Expected {
			if i != 0 {
				label = strings.Repeat(" ", len(label))
			}
			log.Println(yellow(fmt.Sprintf("   %s%s", label, ex)))
		}
		log.Println(green(fmt.Sprintf("     Actual: %s", err.Actual)))

		return
	}
	if err == nil {
		log.Println(red(fmt.Sprintf("Error: %v", tr.Error.Error())))
	} else {
		log.Println(red(fmt.Sprintf("Error: %v", err)))
	}
}

func seqStr(seq int) string {
	return fmt.Sprintf("%d:", seq)
}
