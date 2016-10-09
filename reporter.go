package h2spec

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/summerwind/h2spec/spec"
)

type JUnitTestReport struct {
	XMLName    xml.Name          `xml:"testsuites"`
	TestSuites []*JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	XMLName   xml.Name         `xml:"testsuite"`
	Name      string           `xml:"name,attr"`
	Package   string           `xml:"package,attr"`
	ID        string           `xml:"id,attr"`
	Tests     int              `xml:"tests,attr"`
	Skipped   int              `xml:"skipped,attr"`
	Failures  int              `xml:"failures,attr"`
	Errors    int              `xml:"errors,attr"`
	TestCases []*JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Package   string        `xml:"package,attr"`
	ClassName string        `xml:"classname,attr"`
	Time      string        `xml:"time,attr"`
	Failures  int           `xml:"failures,attr"`
	Failure   *JUnitFailure `xml:"failure"`
	Skipped   *JUnitSkipped `xml:"skipped"`
}

type JUnitFailure struct {
	XMLName xml.Name `xml:"failure"`
	Content string   `xml:",innerxml"`
}

type JUnitSkipped struct {
	XMLName xml.Name `xml:"skipped"`
	Content string   `xml:",innerxml"`
}

type JUnitReporter struct{}

func NewJUnitReporter() *JUnitReporter {
	return &JUnitReporter{}
}

func (e JUnitReporter) Export(results []*spec.TestResult) (string, error) {
	if len(results) == 0 {
		return "", nil
	}

	report := JUnitTestReport{
		TestSuites: make([]*JUnitTestSuite, 0, 20),
	}

	var parent *spec.TestGroup
	var ts *JUnitTestSuite

	for _, r := range results {
		if parent != r.TestCase.Parent {
			parent = r.TestCase.Parent

			ts = &JUnitTestSuite{
				Package:   parent.ID(),
				Name:      fmt.Sprintf("%s. %s", parent.Section, parent.Name),
				ID:        parent.Section,
				Tests:     0,
				Skipped:   0,
				Failures:  0,
				Errors:    0,
				TestCases: make([]*JUnitTestCase, 0, 20),
			}
			report.TestSuites = append(report.TestSuites, ts)
		}

		tc := &JUnitTestCase{
			Package:   r.TestCase.Parent.ID(),
			ClassName: r.TestCase.Desc,
			Time:      fmt.Sprintf("%.04f", r.Duration.Seconds()),
		}

		ts.Tests += 1
		if r.Skipped() {
			ts.Skipped += 1
			tc.Skipped = &JUnitSkipped{}
		} else if r.Failed() {
			ts.Failures += 1

			err := r.Error.(*spec.TestError)
			expected := strings.Join(err.Expected, "\n")
			actual := err.Actual

			tc.Failure = &JUnitFailure{
				Content: fmt.Sprintf("Expect:\n%s\nActual:\n%s", expected, actual),
			}
		}

		ts.TestCases = append(ts.TestCases, tc)
	}

	buf, err := xml.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", xml.Header, buf), nil
}
