/*
Package gisty provides a similar functionality of `gh gist` command.

It aims to provide a simple and easy to use interface to interact with GitHub
Gists in the Go applications.

## Note

In this package, the environment variable "GH_TOKEN" or "GITHUB_TOKEN" must be
set to the "personal access token" of the GitHub API. (`gist` scope required)

For GitHub Enterprise users, the environment variable "GH_ENTERPRISE_TOKEN" or
"GITHUB_ENTERPRISE_TOKEN" must also be set with the GitHub API
"authentication token".
*/
package gisty

import (
	"bytes"

	"github.com/KENOS/go-gisty/gisty/buildinfo"
	"github.com/cli/cli/v2/pkg/cmd/api"
	"github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/cli/cli/v2/pkg/cmd/gist/clone"
	"github.com/cli/cli/v2/pkg/cmd/gist/create"
	"github.com/cli/cli/v2/pkg/cmd/gist/delete"
	"github.com/cli/cli/v2/pkg/cmd/gist/list"
	"github.com/cli/cli/v2/pkg/cmd/gist/view"
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
}

// AltFunc is a set of alternative functions to be used in the commands.
//
// Even though it is mostly used for dependency-injection purposes during testing,
// it can be used to overrride the default behavior of the commands.
type AltFunc struct {
	Clone     func(*clone.CloneOptions) error
	Create    func(*create.CreateOptions) error
	Delete    func(*delete.DeleteOptions) error
	List      func(*list.ListOptions) error
	Read      func(*view.ViewOptions) error
	Stargazer func(*api.ApiOptions) error
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// NewGisty returns a new instance of Gisty.
func NewGisty() *Gisty {
	buildDate := buildinfo.Date
	buildVersion := buildinfo.Version
	ios, stdin, stdout, stderr := iostreams.Test()
	cmdFactory := factory.New(buildVersion)

	cmdFactory.IOStreams = ios

	return &Gisty{
		AltFunctions: AltFunc{
			Clone:     nil,
			Create:    nil,
			Delete:    nil,
			List:      nil,
			Read:      nil,
			Stargazer: nil,
		},
		Factory:      cmdFactory,
		Stdin:        stdin,
		Stdout:       stdout,
		Stderr:       stderr,
		BuildDate:    buildDate,
		BuildVersion: buildVersion,
	}
}

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// SanitizeGistID removes non-alphanumeric characters from gistID.
func SanitizeGistID(gistID string) string {
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
