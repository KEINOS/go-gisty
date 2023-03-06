package gisty

import (
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/gist/shared"
	"github.com/cli/cli/v2/pkg/cmd/gist/view"
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

	cmdList := view.NewCmdView(g.Factory, runView)

	cmdList.SetArgs([]string{gist})
	cmdList.SetIn(g.Stdin)
	cmdList.SetOut(g.Stdout)
	cmdList.SetErr(g.Stderr)

	if err := WrapIfErr(cmdList.Execute(), "failed to read gist"); err != nil {
		return nil, err
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

	cfg, err := opts.Config()
	if err != nil || forceFailReadConf {
		// forceFailReadConf is a dependency-injection to force the failure of
		// reading the configuration. It was implemented so, since opts.Config
		// is a function that returns an object of internal package config.Config
		// we can't mock it directly.
		if err == nil && forceFailReadConf {
			err = NewErr("forced error")
		}

		return nil, WrapIfErr(err, "failed to read option config")
	}

	hostname, _ := cfg.DefaultHost()

	gist, err := sharedGetGist(client, hostname, gistID)
	if err != nil {
		return nil, WrapIfErr(err, "failed to get gist")
	}

	return gist, nil
}
