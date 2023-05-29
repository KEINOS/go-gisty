/*
Package main for an example of how to use the package gisty.

This program prints the comments in a gist on GitHub.

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
	"github.com/MakeNowJust/heredoc"
)

const gistIDDefault = "78cc23f37e55e848905fc4224483763d"

func main() {
	gistID := gistIDDefault

	if len(os.Args) > 1 {
		gistID = gisty.SanitizeGistID(os.Args[1])
	}

	obj := gisty.NewGisty()

	comments, err := obj.Comments(gistID)
	if err != nil {
		log.Fatal(err)
	}

	for index, node := range comments {
		fmt.Println(heredoc.Docf(`
		Comment #%d
			Comment ID: %s
			Author:
				Login Name: %s
				Avatar URL: %s
			Association: %v
			Body (raw): %#v
			Body (HTML): %#v
			Body (text): %#v
			Created At: %s
			Last Edited At: %s
		`,
			index,
			node.ID,
			node.Author.Login,
			node.Author.AvatarURL,
			node.AuthorAssociation,
			node.BodyRaw,
			node.BodyHTML,
			node.BodyText,
			node.CreatedAt,
			node.LastEditedAt,
		))
	}
}
