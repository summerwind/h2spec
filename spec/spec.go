package spec

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
)

var (
	// ErrTimeout is used when the test times out.
	ErrTimeout = errors.New("Timeout")
	// ErrSkipped is used when the test skipped.
	ErrSkipped = errors.New("Skipped")
)

// TestGroup represents a group of test case.
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

// IsRoot returns bool as to whether it is the parent of all groups.
func (tg *TestGroup) IsRoot() bool {
	return tg.Parent == nil
}

// ID returns the unique ID of this group.
func (tg *TestGroup) ID() string {
	if tg.IsRoot() {
		return tg.Key
	}

	return fmt.Sprintf("%s/%s", tg.Key, tg.Section)
}

// Title returns the title of this group.
func (tg *TestGroup) Title() string {
	if tg.IsRoot() {
		return fmt.Sprintf("%s", tg.Name)
	} else {
		return fmt.Sprintf("%s. %s", tg.Section, tg.Name)
	}
}

// Level returns a number. Level is determined by Key and the number
// of "." included in Section.
func (tg *TestGroup) Level() int {
	if tg.IsRoot() {
		return 0
	}

	return strings.Count(tg.Section, ".") + 1
}

// Test runs all the tests included in this group.
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

			tested = true
		}
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

// AddTestGroup registers a group to this group.
func (tg *TestGroup) AddTestGroup(stg *TestGroup) {
	stg.Parent = tg
	if tg.Strict {
		stg.Strict = true
	}
	tg.Groups = append(tg.Groups, stg)
}

// AddTestCase registers a test to this group.
func (tg *TestGroup) AddTestCase(tc *TestCase) {
	tc.Parent = tg
	if tg.Strict {
		tc.Strict = true
		tg.StrictTests = append(tg.StrictTests, tc)
	} else {
		tg.Tests = append(tg.Tests, tc)
	}
}

// TestCase represents a test case.
type TestCase struct {
	Desc        string
	Requirement string
	Strict      bool
	Parent      *TestGroup
	Result      *TestResult
	Run         func(c *config.Config, conn *Conn) error
}

// Test runs itself as a test case.
func (tc *TestCase) Test(c *config.Config, seq int) error {
	if tc.Strict && !c.Strict {
		return nil
	}

	mode := c.RunMode(fmt.Sprintf("%s/%d", tc.Parent.ID(), seq))
	if mode == config.RunModeNone {
		return nil
	}

	if c.DryRun {
		msg := fmt.Sprintf("%s %s", seqStr(seq), tc.Desc)
		log.Println(msg)
		tc.Result = NewTestResult(tc, seq, nil, time.Duration(0), nil)
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

	tr := NewTestResult(tc, seq, err, end.Sub(start), conn.LocalAddr())
	tr.Print()
	tc.Result = tr

	return nil
}

// TestError represents a error result of test case and implements
// type error.
type TestError struct {
	Expected []string
	Actual   string
}

// Error returns a string containing the reason of the error.
func (e TestError) Error() string {
	return fmt.Sprintf("%s\n%s", strings.Join(e.Expected, "\n"), e.Actual)
}

// TestResult represents a result of test case.
type TestResult struct {
	TestCase   *TestCase
	Sequence   int
	Error      error
	Duration   time.Duration
	SourceAddr net.Addr

	Skipped bool
	Failed  bool
}

// NewTestResult returns a TestResult.
func NewTestResult(tc *TestCase, seq int, err error, d time.Duration, addr net.Addr) *TestResult {
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
		TestCase:   tc,
		Sequence:   seq,
		Error:      err,
		Duration:   d,
		SourceAddr: addr,
		Skipped:    skipped,
		Failed:     failed,
	}

	return &tr
}

// Print prints the result of test case.
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

	log.Println(red(fmt.Sprintf("using source address %s", tr.SourceAddr)))
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
