package reporter

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/summerwind/h2spec/spec"
)

// JUnitReport represents the JUnit XML format.
type JUnitTestReport struct {
	XMLName    xml.Name          `xml:"testsuites"`
	TestSuites []*JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents the testsuite element of JUnit XML format.
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

// JUnitTestCase represents the testcase element of JUnit XML format.
type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Package   string        `xml:"package,attr"`
	ClassName string        `xml:"classname,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure"`
	Skipped   *JUnitSkipped `xml:"skipped"`
	Error     *JUnitError   `xml:"error"`
}

// JUnitFailure represents the failure element of JUnit XML format.
type JUnitFailure struct {
	XMLName xml.Name `xml:"failure"`
	Content string   `xml:",innerxml"`
}

// JUnitSkipped represents the skipped element of JUnit XML format.
type JUnitSkipped struct {
	XMLName xml.Name `xml:"skipped"`
	Content string   `xml:",innerxml"`
}

// JUnitSkipped represents the error element of JUnit XML format.
type JUnitError struct {
	XMLName xml.Name `xml:"error"`
	Content string   `xml:",innerxml"`
}

// JUnitReport writes a file which contains the JUnit report generated
// by test result of h2spec.
func JUnitReport(groups []*spec.TestGroup, filePath string) error {
	report := JUnitTestReport{
		TestSuites: convertJUnitReport(groups),
	}

	buf, err := xml.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	body := fmt.Sprintf("%s%s", xml.Header, buf)
	return ioutil.WriteFile(filePath, []byte(body), os.ModePerm)
}

func convertJUnitReport(groups []*spec.TestGroup) []*JUnitTestSuite {
	ts := make([]*JUnitTestSuite, 20)

	for _, tg := range groups {
		tests := append(tg.Tests, tg.StrictTests...)
		if len(tests) == 0 {
			ts = append(ts, convertJUnitReport(tg.Groups)...)
			continue
		}

		jts := &JUnitTestSuite{
			Package:   tg.ID(),
			Name:      fmt.Sprintf("%s. %s", tg.Section, tg.Name),
			ID:        tg.Section,
			Tests:     0,
			Skipped:   0,
			Failures:  0,
			Errors:    0,
			TestCases: make([]*JUnitTestCase, 20),
		}

		for _, tc := range tests {
			if tc.Result == nil {
				continue
			}

			jtc := &JUnitTestCase{
				Package:   tg.ID(),
				ClassName: tc.Desc,
				Time:      fmt.Sprintf("%.04f", tc.Result.Duration.Seconds()),
			}

			jts.Tests += 1
			if tc.Result.Skipped {
				jts.Skipped += 1
				jtc.Skipped = &JUnitSkipped{}
			} else if tc.Result.Failed {
				switch tc.Result.Error.(type) {
				case spec.TestError:
					jts.Failures += 1

					err := tc.Result.Error.(*spec.TestError)
					expected := strings.Join(err.Expected, "\n")
					actual := err.Actual

					jtc.Failure = &JUnitFailure{
						Content: fmt.Sprintf("Expect:\n%s\nActual:\n%s", expected, actual),
					}
				default:
					jts.Errors += 1

					err := tc.Result.Error
					jtc.Error = &JUnitError{
						Content: err.Error(),
					}
				}
			}

			jts.TestCases = append(jts.TestCases, jtc)
		}

		ts = append(ts, jts)
		ts = append(ts, convertJUnitReport(tg.Groups)...)
	}

	return ts
}
