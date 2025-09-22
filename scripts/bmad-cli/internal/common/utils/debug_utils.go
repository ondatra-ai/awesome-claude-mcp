package utils

import (
	"fmt"
	"os"
)

func DebugEnabled() bool {
	return os.Getenv("DEBUG") != ""
}

func DebugLogWithSeparator(header, content string) {
	if DebugEnabled() {
		fmt.Fprintf(os.Stderr, "%s:\n", header)
		fmt.Fprintln(os.Stderr, "--------------------------------")
		fmt.Fprintln(os.Stderr, content)
		fmt.Fprintln(os.Stderr, "--------------------------------")
	}
}

func LogDebugf(format string, args ...interface{}) {
	if DebugEnabled() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
