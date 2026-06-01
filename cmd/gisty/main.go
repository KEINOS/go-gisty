// Command gisty forwards gist operations to the GitHub CLI.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/KEINOS/go-gisty/internal/ghcmd"
)

var (
	execCommandContext           = exec.CommandContext
	exit                         = os.Exit
	stdin              io.Reader = os.Stdin
	stdout             io.Writer = os.Stdout
	stderr             io.Writer = os.Stderr
)

func main() {
	err := run(os.Args[1:], ghcmd.Streams{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		_, printErr := fmt.Fprintln(stderr, err)
		if printErr != nil {
			exit(exitCode(err))

			return
		}

		exit(exitCode(err))
	}
}

func exitCode(err error) int {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

func run(args []string, streams ghcmd.Streams) error {
	err := ghcmd.Run(
		context.Background(),
		execCommandContext,
		streams,
		append([]string{"gist"}, args...)...,
	)
	if err != nil {
		return fmt.Errorf("execute gh gist: %w", err)
	}

	return nil
}
