package gisty

import (
	"context"
	"os/exec"

	"github.com/KEINOS/go-gisty/internal/ghcmd"
)

const commandGist = "gist"

var execCommandContext = exec.CommandContext

func (g *Gisty) runGH(args ...string) error {
	return WrapIfErr(
		ghcmd.Run(context.Background(), execCommandContext, g.streams(), args...),
		"failed to execute gh command",
	)
}
