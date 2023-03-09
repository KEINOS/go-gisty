/*
Package main for an example of how to use the package gisty.

This program prints the list of gists on GitHub.

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

	"github.com/KEINOS/go-gisty/gisty"
)

func main() {
	obj := gisty.NewGisty()

	argsList := gisty.ListArgs{
		Limit:      100,
		OnlyPublic: false,
		OnlySecret: false,
	}

	infos, err := obj.List(argsList)
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range infos {
		fmt.Printf(
			"ID: %s\n\tDescription: %s\n\tNumber of files: %v\n\tUpdatedAt: %s\n\tIsPublic: %v\n",
			info.GistID,
			info.Description,
			info.Files,
			info.UpdatedAt,
			info.IsPublic,
		)
	}
}
