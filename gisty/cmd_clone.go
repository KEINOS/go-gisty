package gisty

import (
	"github.com/cli/cli/v2/pkg/cmd/gist/clone"
	"github.com/pkg/errors"
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

func (g *Gisty) clone(args []string, runF func(*clone.CloneOptions) error) error {
	cmdList := clone.NewCmdClone(g.Factory, runF)

	cmdList.SetArgs(args)
	cmdList.SetIn(g.Stdin)
	cmdList.SetOut(g.Stdout)
	cmdList.SetErr(g.Stderr)

	return errors.Wrap(cmdList.Execute(), "failed to clone gist")
}
