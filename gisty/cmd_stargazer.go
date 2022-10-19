package gisty

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/cli/cli/v2/pkg/cmd/api"
	"github.com/pkg/errors"
)

// Stargazer returns the number of stars in the gist for a given gist ID.
//
// Note that gistID should not be the gist URL.
func (g *Gisty) Stargazer(gistID string) (int, error) {
	return g.stargazer(gistID, g.AltFunctions.Stargazer)
}

// Stargazer returns the number of stars in the gist for a given gist ID.
func (g *Gisty) stargazer(gistID string, runF func(*api.ApiOptions) error) (int, error) {
	cmdAPI := api.NewCmdApi(g.Factory, runF)

	query := fmt.Sprintf(
		"query { viewer { gist (name: \"%s\" ) { name, stargazerCount } } }",
		sanitizeGistID(gistID),
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

	if err := cmdAPI.Execute(); err != nil {
		return 0, errors.Wrap(err, "failed to execute GitHub API request")
	}

	count, err := strconv.Atoi(strings.Trim(g.Stdout.String(), "'"))
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse GitHub API response.\nAPI response=%#v", g.Stdout.String())
		err = errors.Wrap(err, errMsg)
	}

	return count, err
}

// sanitizeGistID removes non-alphanumeric characters from gistID.
func sanitizeGistID(gistID string) string {
	bytesGistID := []byte(gistID)
	index := 0

	for _, b := range bytesGistID {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			bytesGistID[index] = b
			index++
		}
	}

	return string(bytesGistID[:index])
}
