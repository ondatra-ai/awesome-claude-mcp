package prtriage

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// exitCode extracts the exit code from an error.
func exitCode(err error) int {
	if err == nil {
		return 0
	}

	ee := &exec.ExitError{}
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}

	return -1
}

// firstLine returns the first line of a string.
func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}

	return s
}

// truncate truncates a string to the given maximum length.
func truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	return s[:maxLength]
}

// debugEnabled checks if debug mode is enabled via environment variable.
func debugEnabled() bool {
	return os.Getenv("DEBUG") != ""
}

// debugLogWithSeparator logs content to stderr with a header and separators if debug is enabled.
func debugLogWithSeparator(header, content string) {
	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "%s:\n", header)
		fmt.Fprintln(os.Stderr, "--------------------------------")
		fmt.Fprintln(os.Stderr, content)
		fmt.Fprintln(os.Stderr, "--------------------------------")
	}
}

// logDebugf logs a simple debug message to stderr if debug is enabled.
func logDebugf(format string, args ...interface{}) {
	if debugEnabled() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
