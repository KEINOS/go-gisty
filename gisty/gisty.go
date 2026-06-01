/*
Package gisty provides a similar functionality of `gh gist` command.

It aims to provide a simple and easy to use interface to interact with GitHub
Gists in the Go applications.

- Note

In this package, the environment variable "GH_TOKEN" or "GITHUB_TOKEN" must be
set to the "personal access token" of the GitHub API. (`gist` scope required)

For GitHub Enterprise users, the environment variable "GH_ENTERPRISE_TOKEN" or
"GITHUB_ENTERPRISE_TOKEN" must also be set with the GitHub API
"authentication token".

- Tips and info to implement Gisty commands

The basic of Gisty is a wrapper of the `gh gist` command. Gisty receives the
command arguments and passes them to the installed `gh` application. The `Read`
method uses the GitHub API directly because it returns a structured Gist value.

Commands that are not supported by `gh`, `Stargazer` for example, are implemented
via the `gh api` sub command. This sub command supports GitHub GraphQL API (v4),
which supports more features
than the REST API (v3) that `gh` uses.

- GraphQL API documentation: https://docs.github.com/en/graphql
- GraphQL Explorer: https://docs.github.com/ja/graphql/overview/explorer
*/
package gisty

import (
	"bytes"

	buildinfo "github.com/KEINOS/go-gisty/gisty/buildinfos"
	"github.com/KEINOS/go-gisty/gisty/internal/gistid"
	"github.com/KEINOS/go-gisty/gisty/internal/httpclient"
	"github.com/cli/cli/v2/pkg/cmd/api"
	"github.com/cli/cli/v2/pkg/cmd/gist/clone"
	"github.com/cli/cli/v2/pkg/cmd/gist/create"
	"github.com/cli/cli/v2/pkg/cmd/gist/delete"
	"github.com/cli/cli/v2/pkg/cmd/gist/list"
	"github.com/cli/cli/v2/pkg/cmd/gist/view"
	"github.com/cli/cli/v2/pkg/cmd/repo/sync"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
)

// ----------------------------------------------------------------------------
//  Type: Gisty, AltFunc
// ----------------------------------------------------------------------------

// Gisty is the main struct of this package.
type Gisty struct {
	// AltFunctions is a set of alternative functions to be used in the
	// commands. If nil is set, the default function is used.
	AltFunctions AltFunc
	// Factory holds the I/O streams, http client, and other common
	// dependencies to request GitHub API.
	Factory *cmdutil.Factory
	// Stdin is the standard input stream which each command reads from.
	Stdin *bytes.Buffer
	// Stdout is the standard output stream which each command writes to.
	Stdout *bytes.Buffer
	// Stderr is the standard error stream which each command writes to.
	Stderr *bytes.Buffer
	// BuildDate is the date when the binary was built.
	BuildDate string
	// BuildVersion is the version of the binary.
	BuildVersion string
	// MaxComment is the max number of comments in a gist to be fetched.
	MaxComment int
}

// AltFunc is a set of alternative functions to be used in the commands.
//
// Even though it is mostly used for dependency-injection purposes during testing,
// it can be used to overrride the default behavior of the commands.
type AltFunc struct {
	Clone     func(*clone.CloneOptions) error
	Comments  func(*api.ApiOptions) error
	Create    func(*create.CreateOptions) error
	Delete    func(*delete.DeleteOptions) error
	List      func(*list.ListOptions) error
	Read      func(*view.ViewOptions) error
	Stargazer func(*api.ApiOptions) error
	Update    func(*sync.SyncOptions) error
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// MaxCommentDefault is the default value of max number of comments to be fetched.
const MaxCommentDefault = 100

// NewGisty returns a new instance of Gisty.
func NewGisty() *Gisty {
	altFn := new(AltFunc)
	buildDate := buildinfo.Date
	buildVersion := buildinfo.Version
	ios, stdin, stdout, stderr := iostreams.Test()

	cmdFactory := new(cmdutil.Factory)

	cmdFactory.AppVersion = buildVersion
	cmdFactory.InvokingAgent = "go-gisty"
	cmdFactory.IOStreams = ios
	cmdFactory.HttpClient = httpclient.New(buildVersion, "go-gisty")

	gst := new(Gisty)

	gst.AltFunctions = *altFn
	gst.Factory = cmdFactory
	gst.Stdin = stdin
	gst.Stdout = stdout
	gst.Stderr = stderr
	gst.BuildDate = buildDate
	gst.BuildVersion = buildVersion
	gst.MaxComment = MaxCommentDefault

	return gst
}

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// SanitizeGistID removes non-alphanumeric characters from gistID.
func SanitizeGistID(gistID string) string {
	return gistid.Sanitize(gistID)
}
