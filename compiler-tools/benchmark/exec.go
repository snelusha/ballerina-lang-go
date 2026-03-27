package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// execQuiet runs a command in the given directory, capturing stdout and
// stderr. On failure it returns an error that includes the captured output.
func execQuiet(dir, name string, args ...string) error {
	var buf bytes.Buffer
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = &buf
	c.Stderr = &buf
	if err := c.Run(); err != nil {
		return fmt.Errorf("%w\n%s", err, buf.String())
	}
	return nil
}

// execCapture runs a command and returns its combined output as a string.
func execCapture(dir, name string, args ...string) (string, error) {
	var buf bytes.Buffer
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = &buf
	c.Stderr = &buf
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("%w\n%s", err, buf.String())
	}
	return buf.String(), nil
}
