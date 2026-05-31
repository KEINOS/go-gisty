package gisty

import (
	"context"
	"os/exec"
)

const commandGist = "gist"

var execCommandContext = exec.CommandContext

func (g *Gisty) runGH(args ...string) error {
	cmd := execCommandContext(context.Background(), "gh", args...)
	cmd.Stdin = g.Stdin
	cmd.Stdout = g.Stdout
	cmd.Stderr = g.Stderr

	return WrapIfErr(cmd.Run(), "failed to execute gh command")
}
