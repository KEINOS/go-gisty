package gisty

import (
	"github.com/cli/cli/v2/pkg/cmd/gist/delete"
	"github.com/pkg/errors"
)

// Delete deletes a gist for a given gist ID or URL.
//
// Note that it will remove the gist right away, without any confirmation.
func (g *Gisty) Delete(gist string) error {
	return g.delete(gist, g.AltFunctions.Delete)
}

func (g *Gisty) delete(gist string, runF func(*delete.DeleteOptions) error) error {
	cmdList := delete.NewCmdDelete(g.Factory, runF)

	cmdList.SetArgs([]string{gist})
	cmdList.SetIn(g.Stdin)
	cmdList.SetOut(g.Stdout)
	cmdList.SetErr(g.Stderr)

	return errors.Wrap(cmdList.Execute(), "failed to delete gist")
}
