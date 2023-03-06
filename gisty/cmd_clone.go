package gisty

import (
	"github.com/cli/cli/v2/pkg/cmd/gist/clone"
)

// Clone clones a gist with the given args.
//
//	[]string{
//	  <gist>,             // required
//	  [<directory>],      // optional
//	  [-- <gitflags>...], // optional
//	}
//
// <gist> is a gist ID or URL and <directory> is the directory to clone the gist
// into. <gitflags> listed after '--' are the additional flags to pass directly
// as a 'git clone' flags.
func (g *Gisty) Clone(args []string) error {
	return g.clone(args, g.AltFunctions.Clone)
}

// clone is a wrapper around the clone command from the gh cli.
//
// If altF is not nil, it will be used instead of the default function.
func (g *Gisty) clone(args []string, altF func(*clone.CloneOptions) error) error {
	cmdList := clone.NewCmdClone(g.Factory, altF)

	cmdList.SetArgs(args)
	cmdList.SetIn(g.Stdin)
	cmdList.SetOut(g.Stdout)
	cmdList.SetErr(g.Stderr)

	return WrapIfErr(cmdList.Execute(), "failed to execute gist clone")
}
