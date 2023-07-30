package gisty

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/pkg/cmd/api"
)

// Template for the GraphQL query to get the comments in the gist.
const tplQueryComments = `
{
	viewer {
		gist(name: "%s") {
			comments(last: %d) {
				edges {
					node {
						id
						authorAssociation
						author {
							avatarUrl
							login
						}
						createdAt
						publishedAt
						lastEditedAt
						body
						bodyHTML
						bodyText
						isMinimized
						minimizedReason
					}
				}
			}
		}
	}
}`

// ----------------------------------------------------------------------------
//  Dummy data for testing
// ----------------------------------------------------------------------------

// DummyID is the ID of the dummy gist used in the example of the test.
const DummyID = "42f5f23053ab59ca480f480b8d01e1fd"

// DummyComment is the dummy comment used in the example of the test.
//
//nolint:lll // long line is intentional
var DummyComment = Comment{
	Author: Author{
		AvatarURL: "https://avatars.githubusercontent.com/u/11840938?u=e915b35bd36abfdcbbaaa6fbe5ea0c6e8ee51e70&v=4",
		Login:     "KEINOS",
	},
	ID:                "GC_lADOALStqtoAIDQyZjVmMjMwNTNhYjU5Y2E0ODBmNDgwYjhkMDFlMWZkzgBF6l4",
	AuthorAssociation: "OWNER",
	BodyRaw:           "1st example comment @ 20230528.\r\n\r\n- This line was added by edit.",
	BodyHTML:          "<p dir=\"auto\">1st example comment @ 20230528.</p>\n<ul dir=\"auto\">\n<li>This line was added by edit.</li>\n</ul>",
	BodyText:          "1st example comment @ 20230528.\n\nThis line was added by edit.",
	CreatedAt:         "2023-05-28T08:36:32Z",
	PublishedAt:       "2023-05-28T08:36:32Z",
	LastEditedAt:      "2023-05-28T08:44:10Z",
	IsMinimized:       false,
	MinimizedReason:   "",
}

// ----------------------------------------------------------------------------
//  Type: Author
// ----------------------------------------------------------------------------

type Author struct {
	AvatarURL string `json:"avatarUrl"`
	Login     string `json:"login"`
}

// ----------------------------------------------------------------------------
//  Type: Comment
// ----------------------------------------------------------------------------

// Comment is a struct for the comment node in the GraphQL response.
//
//nolint:tagliatelle // bodyHtml is ideal but bodyHTML is required by GitHub API
type Comment struct {
	Author            Author `json:"author"`
	ID                string `json:"id"`
	AuthorAssociation string `json:"authorAssociation"`
	BodyRaw           string `json:"body"`
	BodyHTML          string `json:"bodyHTML"`
	BodyText          string `json:"bodyText"`
	CreatedAt         string `json:"createdAt"`
	PublishedAt       string `json:"publishedAt"`
	LastEditedAt      string `json:"lastEditedAt"`
	MinimizedReason   string `json:"minimizedReason"`
	IsMinimized       bool   `json:"isMinimized"`
}

// ----------------------------------------------------------------------------
/// Methods for the Gisty type
// ----------------------------------------------------------------------------

// Comments returns the comments in the gist.
func (g *Gisty) Comments(gistID string) ([]Comment, error) {
	// Return dummy data if the gist ID is the dummy ID to avoid unwanted request
	// in the example of the test.
	if gistID == DummyID {
		return []Comment{DummyComment}, nil
	}

	return g.comments(gistID, g.AltFunctions.Comments)
}

// comments is the actual function that gets the comments in the gist.
//
// If altF is not nil, it will be used instead of the default function.
func (g *Gisty) comments(gistID string, runF func(*api.ApiOptions) error) ([]Comment, error) {
	gistID = SanitizeGistID(gistID) // sanitize to avoid unwanted query to request
	if gistID == "" {
		return nil, NewErr("invalid gist ID")
	}

	cmdAPI := api.NewCmdApi(g.Factory, runF)
	query := heredoc.Docf(tplQueryComments, gistID, g.MaxComment)

	argv := []string{
		"graphql",
		"-f", "query=" + query,
		"--jq", "[.data.viewer.gist.comments.edges[].node]",
	}

	cmdAPI.SetArgs(argv)
	cmdAPI.SetIn(g.Stdin)
	cmdAPI.SetOut(g.Stdout)
	cmdAPI.SetErr(g.Stderr)

	// Request the GitHub API.
	err := WrapIfErr(cmdAPI.Execute(), "failed to execute GitHub API request")
	if err != nil {
		return nil, err
	}

	response := g.Stdout.Bytes()

	var nodes []Comment

	if err := json.Unmarshal(response, &nodes); err != nil {
		return nil, WrapIfErr(err, "failed to parse GitHub API response. malformed JSON")
	}

	// Parse the GitHub API response to get the number of stars.
	return nodes, nil
}
