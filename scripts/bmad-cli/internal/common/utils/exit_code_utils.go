package utils

import (
	"errors"
	"os/exec"
)

func ExtractExitCode(err error) int {
	if err == nil {
		return 0
	}

	ee := &exec.ExitError{}
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}

	return -1
}
