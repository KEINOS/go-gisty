// Package ghcmd centralizes GitHub CLI command execution.
package ghcmd

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

// Command is the shared subset of Cobra commands used by go-gisty.
type Command interface {
	SetArgs(args []string)
	SetIn(stdin io.Reader)
	SetOut(stdout io.Writer)
	SetErr(stderr io.Writer)
	Execute() error
}

// Executor creates an external command.
type Executor func(context.Context, string, ...string) *exec.Cmd

// Streams contains the standard streams used by a command.
type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Execute configures and runs an in-process GitHub CLI command.
func Execute(cmd Command, args []string, streams Streams) error {
	cmd.SetArgs(args)
	cmd.SetIn(streams.Stdin)
	cmd.SetOut(streams.Stdout)
	cmd.SetErr(streams.Stderr)

	err := cmd.Execute()
	if err != nil {
		return fmt.Errorf("execute GitHub CLI command: %w", err)
	}

	return nil
}

// Run executes the external gh command with the given arguments.
func Run(ctx context.Context, executor Executor, streams Streams, args ...string) error {
	cmd := executor(ctx, "gh", args...)
	cmd.Stdin = streams.Stdin
	cmd.Stdout = streams.Stdout
	cmd.Stderr = streams.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("run gh command: %w", err)
	}

	return nil
}
