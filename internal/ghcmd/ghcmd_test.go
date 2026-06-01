package ghcmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	commandGist = "gist"
	commandList = "list"
)

var errForced = errors.New("forced error")

type stubCommand struct {
	args    []string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	execute func() error
}

func (s *stubCommand) SetArgs(args []string) {
	s.args = args
}

func (s *stubCommand) SetIn(stdin io.Reader) {
	s.stdin = stdin
}

func (s *stubCommand) SetOut(stdout io.Writer) {
	s.stdout = stdout
}

func (s *stubCommand) SetErr(stderr io.Writer) {
	s.stderr = stderr
}

func (s *stubCommand) Execute() error {
	return s.execute()
}

func TestExecute(t *testing.T) {
	t.Parallel()

	streams := Streams{
		Stdin:  bytes.NewBufferString("in"),
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
	}
	cmd := &stubCommand{
		args:   nil,
		stdin:  nil,
		stdout: nil,
		stderr: nil,
		execute: func() error {
			return nil
		},
	}

	err := Execute(cmd, []string{commandGist, commandList}, streams)

	require.NoError(t, err)
	require.Equal(t, []string{commandGist, commandList}, cmd.args)
	require.Same(t, streams.Stdin, cmd.stdin)
	require.Same(t, streams.Stdout, cmd.stdout)
	require.Same(t, streams.Stderr, cmd.stderr)
}

func TestExecute_error(t *testing.T) {
	t.Parallel()

	cmd := &stubCommand{
		args:   nil,
		stdin:  nil,
		stdout: nil,
		stderr: nil,
		execute: func() error {
			return errForced
		},
	}

	err := Execute(cmd, nil, Streams{
		Stdin:  nil,
		Stdout: nil,
		Stderr: nil,
	})

	require.ErrorContains(t, err, "execute GitHub CLI command")
	require.ErrorContains(t, err, "forced error")
}

func TestRun(t *testing.T) {
	t.Parallel()

	stdout := new(bytes.Buffer)
	streams := Streams{
		Stdin:  bytes.NewBuffer(nil),
		Stdout: stdout,
		Stderr: new(bytes.Buffer),
	}

	var (
		gotName string
		gotArgs []string
	)

	executor := func(ctx context.Context, name string, args ...string) *exec.Cmd {
		gotName = name
		gotArgs = args

		return exec.CommandContext(ctx, "go", "version")
	}

	err := Run(context.Background(), executor, streams, commandGist, commandList)

	require.NoError(t, err)
	require.Equal(t, "gh", gotName)
	require.Equal(t, []string{commandGist, commandList}, gotArgs)
	require.Contains(t, stdout.String(), "go version")
}

func TestRun_error(t *testing.T) {
	t.Parallel()

	executor := func(ctx context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "go", "not-a-command")
	}

	err := Run(context.Background(), executor, Streams{
		Stdin:  nil,
		Stdout: nil,
		Stderr: nil,
	})

	require.ErrorContains(t, err, "run gh command")
}
