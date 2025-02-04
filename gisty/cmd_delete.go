package gisty

import (
	"github.com/cli/cli/v2/pkg/cmd/gist/delete"
)

// Current `gh gist delete“ command requires `--yes` arg option when not running
// interactively. It confirms deletion without prompting.
const argOptYes = "--yes"

// Delete deletes a gist for a given gist ID or URL.
//
// Note that it will remove the gist right away, without any confirmation.
func (g *Gisty) Delete(gist string) error {
	return g.delete(gist, g.AltFunctions.Delete)
}

// delete is a wrapper around the delete command from the gh cli.
//
// If altF is not nil, it will be used instead of the default delete function.
func (g *Gisty) delete(gist string, altF func(*delete.DeleteOptions) error) error {
	cmd := delete.NewCmdDelete(g.Factory, altF)

	cmd.SetArgs([]string{
		argOptYes,
		gist,
	})
	cmd.SetIn(g.Stdin)
	cmd.SetOut(g.Stdout)
	cmd.SetErr(g.Stderr)

	return WrapIfErr(cmd.Execute(), "failed to delete gist")
}
