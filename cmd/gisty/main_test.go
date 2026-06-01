package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/KEINOS/go-gisty/internal/ghcmd"
	"github.com/stretchr/testify/require"
)

const (
	commandGisty = "gisty"
	commandList  = "list"
)

var errForcedWrite = errors.New("forced write error")

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) {
	return 0, errForcedWrite
}

//nolint:paralleltest // This test replaces the package-level command executor.
func TestRun(t *testing.T) {
	oldExecCommandContext := execCommandContext

	t.Cleanup(func() {
		execCommandContext = oldExecCommandContext
	})

	var gotArgs []string

	execCommandContext = func(ctx context.Context, _ string, args ...string) *exec.Cmd {
		gotArgs = args

		return exec.CommandContext(ctx, "go", "version")
	}

	err := run([]string{commandList, "--limit=1"}, ghcmd.Streams{
		Stdin:  bytes.NewBuffer(nil),
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
	})

	require.NoError(t, err)
	require.Equal(t, []string{"gist", commandList, "--limit=1"}, gotArgs)
}

//nolint:paralleltest // This test replaces package-level process dependencies.
func TestMain_error(t *testing.T) {
	oldArgs := os.Args
	oldExecCommandContext := execCommandContext
	oldExit := exit
	oldStderr := stderr

	t.Cleanup(func() {
		os.Args = oldArgs
		execCommandContext = oldExecCommandContext
		exit = oldExit
		stderr = oldStderr
	})

	os.Args = []string{commandGisty, commandList}
	execCommandContext = func(ctx context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "go-gisty-command-not-found")
	}

	var exitCode int

	exit = func(code int) {
		exitCode = code
	}

	stderrBuffer := new(bytes.Buffer)
	stderr = stderrBuffer

	main()

	require.Equal(t, 1, exitCode)
	require.Contains(t, stderrBuffer.String(), "executable file not found")
}

//nolint:paralleltest // This test replaces package-level process dependencies.
func TestMain_stderr_error(t *testing.T) {
	oldArgs := os.Args
	oldExecCommandContext := execCommandContext
	oldExit := exit
	oldStderr := stderr

	t.Cleanup(func() {
		os.Args = oldArgs
		execCommandContext = oldExecCommandContext
		exit = oldExit
		stderr = oldStderr
	})

	os.Args = []string{commandGisty, commandList}
	execCommandContext = func(ctx context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "go-gisty-command-not-found")
	}

	var exitCode int

	exit = func(code int) {
		exitCode = code
	}

	stderr = errorWriter{}

	main()

	require.Equal(t, 1, exitCode)
}

//nolint:paralleltest // This test replaces package-level process dependencies.
func TestMain_preserves_child_exit_code(t *testing.T) {
	oldArgs := os.Args
	oldExecCommandContext := execCommandContext
	oldExit := exit
	oldStderr := stderr
	pathTestBinary := os.Args[0]

	t.Cleanup(func() {
		os.Args = oldArgs
		execCommandContext = oldExecCommandContext
		exit = oldExit
		stderr = oldStderr
	})

	os.Args = []string{commandGisty, commandList}
	execCommandContext = func(ctx context.Context, _ string, _ ...string) *exec.Cmd {
		//nolint:gosec // The path is the current controlled test binary.
		return exec.CommandContext(ctx, pathTestBinary, "-test.run=TestHelperProcess", "--", "exit", "42")
	}

	var exitCode int

	exit = func(code int) {
		exitCode = code
	}

	stderr = new(bytes.Buffer)

	main()

	require.Equal(t, 42, exitCode)
}

//nolint:paralleltest // This test is executed as a subprocess.
func TestHelperProcess(t *testing.T) {
	if len(os.Args) < 3 || os.Args[len(os.Args)-2] != "exit" {
		return
	}

	exitCode, err := strconv.Atoi(os.Args[len(os.Args)-1])
	require.NoError(t, err)

	os.Exit(exitCode)
}
