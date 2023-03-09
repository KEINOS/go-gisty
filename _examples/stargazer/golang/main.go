/*
Package main for an example of how to use the package gisty.

This program prints the number of stargazers of a gist on GitHub.

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
	"os"

	"github.com/KEINOS/go-gisty/gisty"
)

const gistIDDefault = "5b10b34f87955dfc86d310cd623a61d1"

func main() {
	gistID := gistIDDefault

	if len(os.Args) > 1 {
		gistID = gisty.SanitizeGistID(os.Args[1])
	}

	obj := gisty.NewGisty()

	count, err := obj.Stargazer(gistID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s, STARS: %v\n", gistID, count)
}
