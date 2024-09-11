package gisty

import (
	"fmt"
	"strconv"
	"strings"

	"al.essio.dev/pkg/shellescape"
	"github.com/cli/cli/v2/pkg/cmd/api"
)

// Stargazer returns the number of stars in the gist for a given gist ID.
//
// Note that gistID should not be the gist URL.
func (g *Gisty) Stargazer(gistID string) (int, error) {
	return g.stargazer(gistID, g.AltFunctions.Stargazer)
}

// stargazer is a wrapper around the api command from the gh cli to request the
// number of stars for a given gist.
//
// If altF is not nil, it will be used instead of the default function.
func (g *Gisty) stargazer(gistID string, runF func(*api.ApiOptions) error) (int, error) {
	cmdAPI := api.NewCmdApi(g.Factory, runF)
	query := fmt.Sprintf(
		"query { viewer { gist (name: \"%s\" ) { name, stargazerCount } } }",
		SanitizeGistID(gistID), // sanitize to avoid unwanted query to request
	)
	template := "{{.data.viewer.gist.stargazerCount}}"
	argv := []string{
		"graphql",
		"-f", "query=" + query,
		"--template=" + shellescape.Quote(template),
	}

	cmdAPI.SetArgs(argv)
	cmdAPI.SetIn(g.Stdin)
	cmdAPI.SetOut(g.Stdout)
	cmdAPI.SetErr(g.Stderr)

	err := WrapIfErr(cmdAPI.Execute(), "failed to execute GitHub API request")
	if err != nil {
		return 0, err
	}

	// Parse the GitHub API response to get the number of stars.
	count, err := strconv.Atoi(strings.Trim(g.Stdout.String(), "'"))

	return count, WrapIfErr(err, "failed to parse GitHub API response.\nAPI response=%#v", g.Stdout.String())
}
