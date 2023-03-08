package gisty

import (
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/repo/sync"
)

// UpdateArgs are the arguments for the Update function.
type UpdateArgs struct {
	PathDirRepo string // path to the local repository to sync.
	Branch      string // if set, syncs local repo from remote parent on specific branch.
	Destination string // If set, syncs as a destination repo. If not set, syncs to remote parent.
	Source      string // if set, syncs remote repo from another remote source repo.
	Force       bool   // if true, syncs using a hard reset. fast forward update if false.
}

func NewUpdateArgs(pathDirRepo string) UpdateArgs {
	return UpdateArgs{
		PathDirRepo: pathDirRepo,
		Branch:      "",
		Destination: "",
		Source:      "",
		Force:       false,
	}
}

// Update syncs the gist repository with the given args. The returned string is
// the output of the command on success.
//
//nolint:nonamedreturns // named retrun is intentional due to the error from defer
func (g *Gisty) Update(args UpdateArgs) (msg string, err error) {
	if args.PathDirRepo == "" {
		return "", NewErr("path to local repository is required")
	}

	// Change the working directory to the local repository.
	returnPath, err := ChDir(args.PathDirRepo)
	if err != nil {
		return "", WrapIfErr(err, "failed to change working directory to %s", args.PathDirRepo)
	}

	defer func() {
		_, errChDir := ChDir(returnPath)
		if err == nil {
			err = WrapIfErr(errChDir, "failed to change working directory back to %s", returnPath)
		}
	}()

	argsUpdate := []string{}

	if args.Branch != "" {
		argsUpdate = append(argsUpdate, "--branch="+args.Branch)
	}

	if args.Destination != "" {
		argsUpdate = append(argsUpdate, args.Destination)
	}

	if args.Source != "" {
		argsUpdate = append(argsUpdate, "--source="+args.Source)
	}

	if args.Force {
		argsUpdate = append(argsUpdate, "--force")
	}

	return g.update(argsUpdate, g.AltFunctions.Update)
}

// update is a wrapper around the repo.sync command from the gh cli.
//
// If altF is not nil, it will be used instead of the default function.
func (g *Gisty) update(args []string, altF func(*sync.SyncOptions) error) (string, error) {
	cmdSync := sync.NewCmdSync(g.Factory, altF)

	cmdSync.SetArgs(args)
	cmdSync.SetIn(g.Stdin)
	cmdSync.SetOut(g.Stdout)
	cmdSync.SetErr(g.Stderr)

	err := WrapIfErr(cmdSync.Execute(), "failed to execute update/sync command")
	if err != nil {
		return "", err
	}

	// Capture the result of the command execution.
	const successMsgPfx = "✓ Synced"

	result := g.Stdout.String()

	if !strings.HasPrefix(result, successMsgPfx) {
		return "", NewErr("failed to sync gist. Output: '%s'", result)
	}

	return result, nil
}
