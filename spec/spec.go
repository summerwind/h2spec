package spec

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/spec/log"
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
	parent      *TestGroup
	groups      []*TestGroup
	tests       []*TestCase
	strictTests []*TestCase
}

func (tg *TestGroup) ID() string {
	if tg.IsRoot() {
		return tg.Key
	}

	return fmt.Sprintf("%s/%s", tg.Key, tg.Section)
}

func (tg *TestGroup) AddTestGroup(stg *TestGroup) {
	stg.parent = tg
	stg.Strict = tg.Strict
	tg.groups = append(tg.groups, stg)
}

func (tg *TestGroup) AddTestCase(tc *TestCase) {
	tc.parent = tg
	if tg.Strict {
		tc.Strict = true
		tg.strictTests = append(tg.strictTests, tc)
	} else {
		tg.tests = append(tg.tests, tc)
	}
}

func (tg *TestGroup) IsTarget(targets map[string]bool) bool {
	if len(targets) == 0 {
		return true
	}

	_, ok := targets[tg.ID()]
	if ok {
		return true
	}

	key := tg.Key
	val, ok := targets[key]
	if ok && val {
		return true
	}

	if !tg.IsRoot() {
		nums := strings.Split(tg.parent.Section, ".")
		for i, _ := range nums {
			id := fmt.Sprintf("%s/%s", key, strings.Join(nums[:i+1], "."))
			val, ok := targets[id]
			if ok && val {
				return true
			}
		}
	}

	return false
}

//func (tg *TestGroup) IsTarget(targets []string) bool {
//	if len(targets) == 0 {
//		return true
//	}
//
//	id := tg.ID()
//	for _, target := range targets {
//		var base, prefix string
//
//		if len(target) > len(id) {
//			base = target
//			prefix = id
//		} else {
//			base = id
//			prefix = target
//		}
//
//		if strings.HasPrefix(base, prefix) {
//			return true
//		}
//	}
//
//	return false
//}

func (tg *TestGroup) IsRoot() bool {
	return tg.parent == nil
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

func (tg *TestGroup) Test(c *config.Config) []*TestResult {
	results := []*TestResult{}
	level := tg.Level()

	if tg.Strict && !c.Strict {
		return nil
	}

	if !tg.IsTarget(c.Targets) {
		return nil
	}

	log.SetIndentLevel(level)
	log.Println(tg.Title())
	log.SetIndentLevel(level + 1)

	tests := append(tg.tests, tg.strictTests...)
	tested := 0

	for i, tc := range tests {
		num := i + 1

		if c.DryRun {
			log.DescDryRun(num, tc.Desc)
			tested += 1
			continue
		}

		if tc.Strict && !c.Strict {
			continue
		}

		if !tc.IsTarget(num, c.Targets) {
			continue
		}

		if !c.Verbose {
			log.DescRunning(num, tc.Desc)
		}

		r, err := tc.Test(c)
		if err != nil {
			log.ResetLine()
			log.DescError(num, tc.Desc, err)
			os.Exit(1)
		}

		log.ResetLine()
		if r.Skipped() {
			log.DescSkipped(num, tc.Desc)
		} else if r.Failed() {
			err, ok := r.Error.(*TestError)
			if ok {
				log.DescFailed(num, tc.Desc, tc.Requirement, err.Expected, err.Actual)
			} else {
				log.DescError(num, tc.Desc, r.Error)
			}
		} else {
			log.DescPassed(num, tc.Desc)
		}

		results = append(results, r)
		tested += 1
	}

	if tested > 0 {
		log.Println("")
	}

	for _, g := range tg.groups {
		results = append(results, g.Test(c)...)
	}

	return results
}

type TestCase struct {
	Desc        string
	Requirement string
	Strict      bool
	Run         func(c *config.Config, conn *Conn) error
	parent      *TestGroup
}

func (tc *TestCase) IsTarget(num int, targets map[string]bool) bool {
	if len(targets) == 0 {
		return true
	}

	id := fmt.Sprintf("%s/%d", tc.parent.ID(), num)
	val, ok := targets[id]
	if ok && val {
		return true
	}

	key := tc.parent.Key
	val, ok = targets[key]
	if ok && val {
		return true
	}

	nums := strings.Split(tc.parent.Section, ".")
	for i, _ := range nums {
		id := fmt.Sprintf("%s/%s", key, strings.Join(nums[:i+1], "."))
		val, ok := targets[id]
		if ok && val {
			return true
		}
	}

	return false
}

func (tc *TestCase) Test(c *config.Config) (*TestResult, error) {
	conn, err := Dial(c)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	start := time.Now()
	err = tc.Run(c, conn)
	end := time.Now()

	d := end.Sub(start)

	return &TestResult{
		TestCase: tc,
		Error:    err,
		Duration: d,
	}, nil
}

type TestResult struct {
	TestCase *TestCase
	Error    error
	Duration time.Duration
}

func (tr *TestResult) Skipped() bool {
	return tr.Error != nil && tr.Error == ErrSkipped
}

func (tr *TestResult) Failed() bool {
	return tr.Error != nil && tr.Error != ErrSkipped
}

type TestError struct {
	Expected []string
	Actual   string
}

func (e TestError) Error() string {
	return fmt.Sprintf("%s\n%s", strings.Join(e.Expected, "\n"), e.Actual)
}
