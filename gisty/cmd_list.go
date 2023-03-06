package gisty

import (
	"fmt"
	"strings"

	"github.com/cli/cli/v2/pkg/cmd/gist/list"
)

// ListArgs are the arguments/options to the List function.
type ListArgs struct {
	Limit      int  // Limit is the maximum number of gists to fetch (default 10)
	OnlyPublic bool // Show only public gists. Ignored if OnlySecret is true.
	OnlySecret bool // Show only secret gists. Prior than OnlyPublic.
}

// List returns a list of GistInfo objects.
// The returned list depends on the arguments passed to the function.
func (g *Gisty) List(args ListArgs) ([]GistInfo, error) {
	var argsList []string

	if args.Limit > 0 {
		argsList = append(argsList, fmt.Sprintf("--limit=%d", args.Limit))
	}

	if args.OnlyPublic {
		argsList = append(argsList, "--public")
	}

	// Due to the original design, if both OnlyPublic and OnlySecret are true,
	// OnlySecret is prior than OnlyPublic.
	if args.OnlySecret {
		argsList = append(argsList, "--secret")
	}

	return g.list(argsList, g.AltFunctions.List)
}

// list is a wrapper around the list command from the gh cli.
//
// If altF is not nil, it will be used instead of the default delete function.
func (g *Gisty) list(args []string, altF func(*list.ListOptions) error) ([]GistInfo, error) {
	cmdList := list.NewCmdList(g.Factory, altF)

	cmdList.SetArgs(args)
	cmdList.SetIn(g.Stdin)
	cmdList.SetOut(g.Stdout)
	cmdList.SetErr(g.Stderr)

	err := WrapIfErr(cmdList.Execute(), "failed to execute 'gist list' command")
	if err != nil {
		return nil, err
	}

	return parseGistInfo(g.Stdout.String())
}

func parseGistInfo(list string) ([]GistInfo, error) {
	if list == "" {
		return nil, nil
	}

	result := []GistInfo{}

	lines := strings.Split(list, "\n")
	for _, line := range lines {
		if line != "" {
			gistInfo, err := NewGistInfo(line)
			if err != nil {
				return nil, WrapIfErr(err, "failed to parse gist info from: %#v", list)
			}

			result = append(result, gistInfo)
		}
	}

	return result, nil
}
