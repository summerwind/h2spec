package log

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	gray   = color.New(color.FgHiBlack).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

var (
	currentLevel  int    = 0
	currentIndent string = ""
)

func SetIndentLevel(level int) {
	currentLevel = level
	currentIndent = strings.Repeat("  ", level)
}

func Print(a ...interface{}) {
	fmt.Printf("%s%s", currentIndent, fmt.Sprint(a...))
}

func Println(a ...interface{}) {
	fmt.Printf("%s%s", currentIndent, fmt.Sprintln(a...))
}

func DescDryRun(num int, desc string) {
	Println(fmt.Sprintf("%s %s", number(num), desc))
}

func DescRunning(num int, desc string) {
	Print(fmt.Sprintf("  %s %s", gray(number(num)), gray(desc)))
}

func DescPassed(num int, desc string) {
	Println(fmt.Sprintf("%s %s %s", green("✔"), gray(number(num)), gray(desc)))
}

func DescFailed(num int, desc string, req string, expected []string, actual string) {
	Println(fmt.Sprintf("%s %s %s", red("×"), red(number(num)), red(desc)))

	level := currentLevel
	SetIndentLevel(level + 1)
	defer func() {
		SetIndentLevel(level)
	}()

	Println(red(fmt.Sprintf("-> %s", req)))

	label := "Expected: "
	for i, ex := range expected {
		if i != 0 {
			label = strings.Repeat(" ", len(label))
		}
		Println(yellow(fmt.Sprintf("   %s%s", label, ex)))
	}

	Println(green(fmt.Sprintf("     Actual: %s", actual)))
}

func DescSkipped(num int, desc string) {
	Println(fmt.Sprintf("%s %s", cyan(number(num)), cyan(desc)))
}

func DescError(num int, desc string, err error) {
	Println(fmt.Sprintf("%s %s %s", red("×"), red(number(num)), red(desc)))
	Error(err)
}

func Info(msg string) {
	Println(gray(msg))
}

func Verbose(msg string) {
	Println(gray(fmt.Sprintf("   | %s", msg)))
}

func Error(err error) {
	Println(red(fmt.Sprintf("Error: %v", err)))
}

func ResetLine() {
	fmt.Printf("\r")
}

func leftPad(str, padStr string, padLen int) string {
	return strings.Repeat(padStr, padLen-len(str)) + str
}

func number(num int) string {
	return fmt.Sprintf("%d:", num)
}
