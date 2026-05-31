package gisty

import (
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/gist/shared"
	"github.com/cli/cli/v2/pkg/cmd/gist/view"
	ghauth "github.com/cli/go-gh/v2/pkg/auth"
)

// Read returns a list of GistInfo objects. The returned list depends on the
// arguments passed to the function.
func (g *Gisty) Read(gist string) (*shared.Gist, error) {
	return g.read(gist, g.AltFunctions.Read)
}

// read is a wrapper around the read command from the gh cli.
//
// If altF is not nil, it will be used instead of the default function.
func (g *Gisty) read(gist string, altF func(*view.ViewOptions) error) (*shared.Gist, error) {
	resultGist := new(shared.Gist)

	//nolint:nonamedreturns // Named return is intentional.
	runView := func(opts *view.ViewOptions) (err error) {
		resultGist, err = readRun(opts)
		if err != nil {
			return WrapIfErr(err, "failed to execute readRun function")
		}

		return nil
	}

	if altF != nil {
		runView = altF
	}

	cmd := view.NewCmdView(g.Factory, runView)

	cmd.SetArgs([]string{gist})
	cmd.SetIn(g.Stdin)
	cmd.SetOut(g.Stdout)
	cmd.SetErr(g.Stderr)

	err := cmd.Execute()
	if err != nil {
		return nil, WrapIfErr(err, "failed to read gist")
	}

	return resultGist, nil
}

// sharedGetGist is a copy of shared.GetGist to ease testing.
var sharedGetGist = shared.GetGist

// forceFailReadConf is a dependency-injection to force the failure of reading
// the configuration.
var forceFailReadConf = false

func readRun(opts *view.ViewOptions) (*shared.Gist, error) {
	gistID := opts.Selector

	if strings.Contains(gistID, "/") {
		id, err := shared.GistIDFromURL(gistID)
		if err != nil {
			return nil, WrapIfErr(err, "failed to parse gist ID from URL")
		}

		gistID = id
	}

	if gistID == "" {
		return nil, NewErr("no gist specified")
	}

	client, err := opts.HttpClient()
	if err != nil {
		return nil, WrapIfErr(err, "failed to create http client")
	}

	if forceFailReadConf {
		return nil, WrapIfErr(NewErr("forced error"), "failed to read option config")
	}

	hostname, _ := ghauth.DefaultHost()

	gist, err := sharedGetGist(client, hostname, gistID)
	if err != nil {
		return nil, WrapIfErr(err, "failed to get gist")
	}

	return gist, nil
}
