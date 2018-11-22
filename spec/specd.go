package spec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
)

// ClientTestGroup represents a group of test case.
type ClientTestGroup struct {
	Key     string
	Section string
	Name    string
	Parent  *ClientTestGroup
	Groups  []*ClientTestGroup
	Tests   []*ClientTestCase

	PassedCount  int
	FailedCount  int
	SkippedCount int
}

// IsRoot returns bool as to whether it is the parent of all groups.
func (tg *ClientTestGroup) IsRoot() bool {
	return tg.Parent == nil
}

// ID returns the unique ID of this group.
func (tg *ClientTestGroup) ID() string {
	if tg.IsRoot() {
		return tg.Key
	}

	return fmt.Sprintf("%s/%s", tg.Key, tg.Section)
}

// Title returns the title of this group.
func (tg *ClientTestGroup) Title() string {
	if tg.IsRoot() {
		return fmt.Sprintf("%s", tg.Name)
	} else {
		return fmt.Sprintf("%s. %s", tg.Section, tg.Name)
	}
}

// Level returns a number. Level is determined by Key and the number
// of "." included in Section.
func (tg *ClientTestGroup) Level() int {
	if tg.IsRoot() {
		return 0
	}

	return strings.Count(tg.Section, ".") + 1
}

// Test runs all the tests included in this group.
func (tg *ClientTestGroup) Test(c *config.Config) {
	mode := c.RunMode(tg.ID())
	if mode == config.RunModeNone {
		return
	}

	level := tg.Level()

	log.SetIndentLevel(level)
	log.Println(tg.Title())
	log.SetIndentLevel(level + 1)

	for _, tc := range tg.Tests {
		err := tc.Test(c)
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			os.Exit(1)
		}

		if tc.Result == nil {
			// No TestResult found, means the server cannot
			// receive the first request
			fmt.Println("\nError: the server didn't receive the request")
			os.Exit(1)
		}
	}

	for _, g := range tg.Groups {
		g.Test(c)
	}

	log.PrintBlankLine()
}

// AddClientTestGroup registers a group to this group.
func (tg *ClientTestGroup) AddTestGroup(stg *ClientTestGroup) {
	stg.Parent = tg
	tg.Groups = append(tg.Groups, stg)
}

// AddClientTestGroup registers a test to this group.
func (tg *ClientTestGroup) AddTestCase(tc *ClientTestCase) {
	tc.Parent = tg
	tc.Seq = len(tg.Tests) + 1
	tg.Tests = append(tg.Tests, tc)
}

func (tg *ClientTestGroup) ClientTestCases(testCases map[int]*ClientTestCase, c *config.Config, currentPort int) int {
	mode := c.RunMode(tg.ID())
	if mode == config.RunModeNone {
		return currentPort
	}

	for _, tc := range tg.Tests {
		tc.Port = currentPort
		testCases[currentPort] = tc
		currentPort += 1
	}

	for _, g := range tg.Groups {
		currentPort = g.ClientTestCases(testCases, c, currentPort)
	}

	return currentPort
}

func (tg *ClientTestGroup) IncRecursive(failed bool, skipped bool, inc int) {
	if failed {
		tg.FailedCount += inc
	} else if skipped {
		tg.SkippedCount += inc
	} else {
		tg.PassedCount += inc
	}

	if tg.Parent != nil {
		tg.Parent.IncRecursive(failed, skipped, inc)
	}
}

// ClientTestCase represents a test case.
type ClientTestCase struct {
	Seq         int
	Desc        string
	Requirement string
	Parent      *ClientTestGroup
	Result      *ClientTestResult
	Run         func(c *config.Config, conn *Conn) error

	Port int
	Done chan bool
}

// Test runs itself as a test case.
func (tc *ClientTestCase) Test(c *config.Config) error {
	tc.Done = make(chan bool)
	done := make(chan error)
	go func() {
		split := strings.Split(c.Exec, " ")

		binary := split[0]
		args := append(split[1:], tc.FullPath(c))

		cmd := exec.Command(binary, args...)
		if c.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		err := cmd.Run()
		<-tc.Done
		done <- err
	}()

	select {
	case <-done:
		// command failed with non-zero exit code is accept
		if tc.Result != nil {
			log.ResetLine()
			tc.Result.Print()
		}
		return nil
	case <-time.After(time.Duration(3) * time.Second):
		return ErrTimeout
	}

	return nil
}

func (tc *ClientTestCase) FullPath(c *config.Config) string {
	return fmt.Sprintf("%s://%s:%d/", c.Scheme(), c.Host, tc.Port)
}

// ClientTestResult represents a result of test case.
type ClientTestResult struct {
	ClientTestCase *ClientTestCase
	Error          error
	Duration       time.Duration

	Skipped bool
	Failed  bool
}

// NewClientTestResult returns a ClientTestResult.
func NewClientTestResult(tc *ClientTestCase, err error, d time.Duration) *ClientTestResult {
	skipped := false
	failed := false

	if err != nil {
		if err == ErrSkipped {
			skipped = true
		} else {
			failed = true
		}
	}

	tr := ClientTestResult{
		ClientTestCase: tc,
		Error:          err,
		Duration:       d,
		Skipped:        skipped,
		Failed:         failed,
	}

	return &tr
}

// Print prints the result of test case.
func (tr *ClientTestResult) Print() {
	tc := tr.ClientTestCase
	desc := tc.Desc
	seq := seqStr(tc.Seq)

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
