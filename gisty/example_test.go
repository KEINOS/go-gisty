package gisty_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/KEINOS/go-gisty/gisty"
	"github.com/cli/cli/v2/pkg/cmd/api"
	"github.com/cli/cli/v2/pkg/cmd/gist/list"
)

// ============================================================================
//  Examples
// ============================================================================

func ExampleGisty_Comments() {
	// Instantiate a new Gisty object.
	obj := gisty.NewGisty()

	// Get the comments in the gist.
	const gistID = "42f5f23053ab59ca480f480b8d01e1fd"

	comments, err := obj.Comments(gistID)
	if err != nil {
		log.Fatal(err)
	}

	if len(comments) == 0 {
		log.Fatal("No comments found.")
	}

	firstComment := comments[0]

	// Print available fields in the Comment struct.
	for _, field := range []struct {
		nameField string
		value     interface{}
	}{
		// List of fields in the Comment struct.
		{"Name", firstComment.Author.Login},
		{"Icon", firstComment.Author.AvatarURL},
		{"Association", firstComment.AuthorAssociation},
		{"Comment ID", firstComment.ID},
		{"Comment body(Raw)", firstComment.BodyRaw},   // Note that raw body contains "\r\n" for line breaks.
		{"Comment body(HTML)", firstComment.BodyHTML}, // Note that html body contains "\n" for line breaks.
		{"Comment body(Text)", firstComment.BodyText}, // Note that text body contains "\n" for line breaks.
		{"Created at", firstComment.CreatedAt},
		{"Published at", firstComment.PublishedAt},
		{"Updated at", firstComment.LastEditedAt},
		{"Is minimized", firstComment.IsMinimized},
		{"Minimized reason", firstComment.MinimizedReason},
	} {
		val, ok := field.value.(string)
		if ok {
			fmt.Printf("%s: %#v\n", field.nameField, strings.TrimSpace(val))
		} else {
			fmt.Printf("%s: %v\n", field.nameField, field.value)
		}
	}
	// Output:
	// Name: "KEINOS"
	// Icon: "https://avatars.githubusercontent.com/u/11840938?u=e915b35bd36abfdcbbaaa6fbe5ea0c6e8ee51e70&v=4"
	// Association: "OWNER"
	// Comment ID: "GC_lADOALStqtoAIDQyZjVmMjMwNTNhYjU5Y2E0ODBmNDgwYjhkMDFlMWZkzgBF6l4"
	// Comment body(Raw): "1st example comment @ 20230528.\r\n\r\n- This line was added by edit."
	// Comment body(HTML): "<p dir=\"auto\">1st example comment @ 20230528.</p>\n<ul dir=\"auto\">\n<li>This line was added by edit.</li>\n</ul>"
	// Comment body(Text): "1st example comment @ 20230528.\n\nThis line was added by edit."
	// Created at: "2023-05-28T08:36:32Z"
	// Published at: "2023-05-28T08:36:32Z"
	// Updated at: "2023-05-28T08:44:10Z"
	// Is minimized: false
	// Minimized reason: ""
}

// Example to retrieve the list of gists but WITHOUT calling the actual GitHub
// API.
func ExampleGisty_List_dummy_api_call() {
	obj := gisty.NewGisty()

	// Dummy function to avoid calling the actual GitHub API during test/example.
	// Usually, you do not need to set this.
	obj.AltFunctions.List = func(*list.ListOptions) error {
		// Mock the GitHub API response.
		fmt.Fprint(
			obj.Stdout,
			"d5b9800c636dd78defa4f15894d54d29	Title of gist item2	6 files	secret	2022-04-16T06:08:46Z",
		)

		return nil
	}

	// Retrieve the list of gists.
	gistInfos, err := obj.List(gisty.ListArgs{
		Limit:      1000,  // Maximum number of gists to be obtained.
		OnlyPublic: true,  // Get only public gists.
		OnlySecret: false, // Get only secret gists. If true, then prior than OnlyPublic.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the obtained gist information. In this example, only one
	// gist is obtained.
	for numItem, gistInfo := range gistInfos {
		fmt.Printf("#%d GistID: %v\n", numItem+1, gistInfo.GistID)
		fmt.Printf("#%d Description: %v\n", numItem+1, gistInfo.Description)
		fmt.Printf("#%d Num files in a gist: %v\n", numItem+1, gistInfo.Files)
		fmt.Printf("#%d IsPublic: %v\n", numItem+1, gistInfo.IsPublic)
		fmt.Printf("#%d UpdatedAt: %v\n", numItem+1, gistInfo.UpdatedAt)
	}
	// Output:
	// #1 GistID: d5b9800c636dd78defa4f15894d54d29
	// #1 Description: Title of gist item2
	// #1 Num files in a gist: 6
	// #1 IsPublic: false
	// #1 UpdatedAt: 2022-04-16 06:08:46 +0000 UTC
}

// Example to get the number of stars in the gist.
func ExampleGisty_Stargazer() {
	obj := gisty.NewGisty()

	// Dummy function to avoid calling the actual GitHub API during test/example.
	// Usually, you do not need to set this.
	obj.AltFunctions.Stargazer = func(*api.ApiOptions) error {
		numStarsDummy := 10

		// Mock the GitHub API response.
		fmt.Fprintf(obj.Stdout, "'%d'", numStarsDummy)

		return nil
	}

	// Target gist ID to obtain the number of stars.
	gistID := "5b10b34f87955dfc86d310cd623a61d1"

	count, err := obj.Stargazer(gistID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count)
	// Output: 10
}

// ============================================================================
//  Public functions
// ============================================================================

// ----------------------------------------------------------------------------
//  NewErr
// ----------------------------------------------------------------------------

// See: new_err_test.go

// ----------------------------------------------------------------------------
//  NewGistInfo
// ----------------------------------------------------------------------------

func ExampleNewGistInfo() {
	// "line" is one of the lines returned by the `gist list` command as an
	// example. Each line is tab separated and must contain the following
	// information in that order:
	//   1. Gist ID
	//   2. Description
	//   3. Number of files
	//   4. Public or private
	//   5. Updated at (RFC3339/ISO-8601 format)
	line := "7101f542be23e5048198e2a27c3cfda8	Title of gist item1	1 file	public	2022-09-18T18:56:10Z"

	// NewGistInfo parses the line and returns a GistInfo object.
	item, err := gisty.NewGistInfo(line)
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed information.
	fmt.Printf("%T: %v\n", item.GistID, item.GistID)
	fmt.Printf("%T: %v\n", item.Description, item.Description)
	fmt.Printf("%T: %v\n", item.IsPublic, item.IsPublic)
	fmt.Printf("%T: %v\n", item.UpdatedAt, item.UpdatedAt)
	// Output:
	// string: 7101f542be23e5048198e2a27c3cfda8
	// string: Title of gist item1
	// bool: true
	// time.Time: 2022-09-18 18:56:10 +0000 UTC
}

// ----------------------------------------------------------------------------
//  WrapIfErr
// ----------------------------------------------------------------------------

// See: wrap_if_err_test.go
