package gisty

import (
	"net/url"
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/gist/create"
	"github.com/pkg/errors"
)

type CreateArgs struct {
	// Description for this gist
	Description string
	// FilePaths to contain in this gist
	FilePaths []string
	// AsPublic indicates whether this gist should be public or secret.
	// By default, it is secret.
	AsPublic bool
}

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

func (g *Gisty) create(args []string, runF func(*create.CreateOptions) error) (gistURL *url.URL, err error) {
	cmdList := create.NewCmdCreate(g.Factory, runF)

	cmdList.SetArgs(args)
	cmdList.SetIn(g.Stdin)
	cmdList.SetOut(g.Stdout)
	cmdList.SetErr(g.Stderr)

	err = errors.Wrap(cmdList.Execute(), "failed to create gist")
	if err == nil {
		gistURL, err = url.Parse(strings.TrimSpace(g.Stdout.String()))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse gist URL")
		}
	}

	return gistURL, err
}
