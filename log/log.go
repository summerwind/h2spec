package log

import (
	"fmt"
	"strings"
)

var (
	IndentLevel int    = 0
	Indent      string = ""
)

func SetIndentLevel(level int) {
	IndentLevel = level
	Indent = strings.Repeat("  ", level)
}

func Print(a ...interface{}) {
	fmt.Printf("%s%s", Indent, fmt.Sprint(a...))
}

func Println(a ...interface{}) {
	fmt.Printf("%s%s", Indent, fmt.Sprintln(a...))
}

func PrintBlankLine() {
	fmt.Println("")
}

func ResetLine() {
	fmt.Printf("\r")
}
