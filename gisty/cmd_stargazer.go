package gisty

import (
	"fmt"
	"strconv"
	"strings"

	"al.essio.dev/pkg/shellescape"
	"github.com/KEINOS/go-gisty/internal/ghcmd"
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
	gistID = SanitizeGistID(gistID)
	query := fmt.Sprintf(
		"query { viewer { gist (name: \"%s\" ) { name, stargazerCount } } }",
		gistID, // sanitize to avoid unwanted query to request
	)
	template := "{{.data.viewer.gist.stargazerCount}}"

	argv := []string{
		"graphql",
		"-f", "query=" + query,
		"--template=" + shellescape.Quote(template),
	}
	if runF == nil {
		err := WrapIfErr(g.runGH(append([]string{"api"}, argv...)...), "failed to execute GitHub API request")
		if err != nil {
			return 0, err
		}
	} else {
		cmdAPI := api.NewCmdApi(g.Factory, runF)

		err := WrapIfErr(ghcmd.Execute(cmdAPI, argv, g.streams()), "failed to execute GitHub API request")
		if err != nil {
			return 0, err
		}
	}

	// Parse the GitHub API response to get the number of stars.
	count, err := strconv.Atoi(strings.Trim(g.Stdout.String(), "'"))

	return count, WrapIfErr(err, "failed to parse GitHub API response.\nAPI response=%#v", g.Stdout.String())
}
