package log

import (
	"fmt"
	"strings"
)

var (
	// IndentLevel is the value of current indent level.
	IndentLevel int = 0
	// Indent is the string of current indent level.
	Indent string = ""
)

// SetIndentLevel sets the current indent level by integer.
func SetIndentLevel(level int) {
	IndentLevel = level
	Indent = strings.Repeat("  ", level)
}

// Print writes the specified string with indent.
func Print(a ...interface{}) {
	fmt.Printf("%s%s", Indent, fmt.Sprint(a...))
}

// Println writes the specified string. Indent is added and a newline
// is appended.
func Println(a ...interface{}) {
	fmt.Printf("%s%s", Indent, fmt.Sprintln(a...))
}

// PrintBlankLine writes empty line.
func PrintBlankLine() {
	fmt.Println("")
}

// Resetline cancels the previous line.
func ResetLine() {
	fmt.Printf("\r")
}
