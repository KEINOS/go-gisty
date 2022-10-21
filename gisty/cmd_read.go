package gisty

import (
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/gist/shared"
	"github.com/cli/cli/v2/pkg/cmd/gist/view"
	"github.com/pkg/errors"
)

// Read returns a list of GistInfo objects. The returned list depends on the
// arguments passed to the function.
func (g *Gisty) Read(gist string) (*shared.Gist, error) {
	return g.read(gist, g.AltFunctions.Read)
}

func (g *Gisty) read(gist string, runF func(*view.ViewOptions) error) (*shared.Gist, error) {
	resultGist := &shared.Gist{}

	runView := func(opts *view.ViewOptions) (err error) {
		resultGist, err = readRun(opts)
		if err != nil {
			return errors.Wrap(err, "failed to execute readRun function")
		}

		return nil
	}

	if runF != nil {
		runView = runF
	}

	cmdList := view.NewCmdView(g.Factory, runView)

	cmdList.SetArgs([]string{gist})
	cmdList.SetIn(g.Stdin)
	cmdList.SetOut(g.Stdout)
	cmdList.SetErr(g.Stderr)

	if err := cmdList.Execute(); err != nil {
		return nil, errors.Wrap(err, "failed to read gist")
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
			return nil, errors.Wrap(err, "failed to parse gist ID from URL")
		}

		gistID = id
	}

	if gistID == "" {
		return nil, errors.New("no gist specified")
	}

	client, err := opts.HttpClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http client")
	}

	// forceFailReadConf is a dependency-injection to force the failure of reading
	// the configuration. Since opts.Config is a function that returns an object of
	// internal package config.Config we can't mock it directly.
	cfg, err := opts.Config()
	if err != nil || forceFailReadConf {
		if err == nil && forceFailReadConf {
			err = errors.New("forced error")
		}

		return nil, errors.Wrap(err, "failed to read option config")
	}

	hostname, _ := cfg.DefaultHost()

	gist, err := sharedGetGist(client, hostname, gistID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get gist")
	}

	return gist, nil
}
