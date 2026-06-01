package gisty

import (
	"net/url"
	"strings"

	"github.com/KEINOS/go-gisty/internal/ghcmd"
	"github.com/cli/cli/v2/pkg/cmd/gist/create"
)

// CreateArgs are the arguments for the Create function.
type CreateArgs struct {
	// Description for this gist
	Description string
	// FilePaths to contain in this gist
	FilePaths []string
	// AsPublic indicates whether this gist should be public or secret.
	// By default, it is secret.
	AsPublic bool
}

// Create creates a new gist with the given args and returns the URL of the gist.
func (g *Gisty) Create(args CreateArgs) (*url.URL, error) {
	argsCreate := []string{}

	if args.AsPublic {
		argsCreate = append(argsCreate, "--public")
	}

	if args.Description != "" {
		argsCreate = append(argsCreate, "--desc="+args.Description)
	}

	argsCreate = append(argsCreate, args.FilePaths...)

	return g.create(argsCreate, g.AltFunctions.Create)
}

// create is a wrapper around the create command from the gh cli.
//
// If altF is not nil, it will be used instead of the default function.
func (g *Gisty) create(args []string, altF func(*create.CreateOptions) error) (*url.URL, error) {
	if altF == nil {
		err := WrapIfErr(g.runGH(append([]string{commandGist, "create"}, args...)...), "failed to execute create command")
		if err != nil {
			return nil, err
		}

		gistURL, err := url.Parse(strings.TrimSpace(g.Stdout.String()))

		return gistURL, WrapIfErr(err, "failed to parse gist URL")
	}

	cmd := create.NewCmdCreate(g.Factory, altF)

	err := WrapIfErr(ghcmd.Execute(cmd, args, g.streams()), "failed to execute create command")
	if err != nil {
		return nil, err
	}

	// Capture the result of the command execution and parse it.
	gistURL, err := url.Parse(strings.TrimSpace(g.Stdout.String()))

	return gistURL, WrapIfErr(err, "failed to parse gist URL")
}
