package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	stdout, stderr, exitCode := testRunCommand(t, nil, "--version")
	t.Logf("stdout: %s", stdout.String())
	t.Logf("stderr: %s", stderr.String())
	t.Logf("exitCode: %d", exitCode)
	assert.Equal(t, stdout.String(), version+"\n")
	assert.Equal(t, stderr.String(), "")
	assert.Equal(t, exitCode, 0)
}

func testRunCommand(t *testing.T, stdin io.Reader, args ...string) (stdout, stderr *bytes.Buffer, exitCode int) {
	t.Helper()

	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	args = append([]string{"run", "main.go"}, args...)

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	assert.Nil(t, cmd.Run(), fmt.Sprintf("Run go %v", args))
	exitCode = cmd.ProcessState.ExitCode()
	return stdout, stderr, exitCode
}
