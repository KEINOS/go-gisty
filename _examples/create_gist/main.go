/*
Package main for an example of how to use the package gisty.

This program creates a new gist on GitHub and returns the URL of the gist.

## Note

In this package, the **environment variable "GH_TOKEN" or "GITHUB_TOKEN" must be
set** to the "personal access token" of the GitHub API (with `gist` scope is required).

For GitHub Enterprise users, the environment variable "GH_ENTERPRISE_TOKEN" or
"GITHUB_ENTERPRISE_TOKEN" must also be set with the GitHub API "authentication token".
*/
package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/KEINOS/go-gisty/gisty"
)

func main() {
	obj := gisty.NewGisty()

	argsCreate := gisty.CreateArgs{
		Description: "sample description",
		FilePaths: []string{
			filepath.Join("testdata", "foo.md"),
			filepath.Join("testdata", "bar.md"),
		},
		AsPublic: false,
	}

	gistURL, err := obj.Create(argsCreate)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gistURL.Path)
	fmt.Println("OK")

	// Output: OK
}
