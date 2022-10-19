package gisty_test

import (
	"fmt"
	"log"
	"time"

	"github.com/KENOS/go-gisty/gisty"
)

// ============================================================================
//  Examples (E2E tests)
// ============================================================================

// Example of listing gists.
//
// In this example, we are actually requesting the GitHub API (E2E test). This
// is not the recommended method, as it is onerous on the API side, but it is
// included for clarity and ease of use.
func Example() {
	obj := gisty.NewGisty()

	// The below line is equivalent to:
	//   gh gist list --public --limit 10
	args := gisty.ListArgs{
		Limit:      10,
		OnlyPublic: true, // If both OnlyPublic and OnlySecret are true, OnlySecret is prior.
		OnlySecret: false,
	}

	items, err := obj.List(args)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(items))

	// Sleep for 1 second to avoid too many requests to GitHub API
	time.Sleep(time.Second)

	// Output: 10
}

// Example to get the number of stars in the gist.
//
// In this example, we are actually requesting the GitHub API (E2E test). This
// is not the recommended method, as it is onerous on the API side, but it is
// included for clarity and ease of use.
func ExampleGisty_Stargazer() {
	gistID := "5b10b34f87955dfc86d310cd623a61d1"
	obj := gisty.NewGisty()

	count, err := obj.Stargazer(gistID)
	if err != nil {
		log.Fatal(err)
	}

	if count > 1 {
		fmt.Println("OK")
	}

	// Output: OK
}

// ============================================================================
//  Public functions
// ============================================================================

// ----------------------------------------------------------------------------
//  NewGistInfo
// ----------------------------------------------------------------------------

func ExampleNewGistInfo() {
	// This line is one of the lines returned by the `gist list` command.
	// Each line is tab separated and contains the following information:
	//   1. Gist ID
	//   2. Description
	//   3. Number of files
	//   4. Public or private
	//   5. Updated at (RFC3339/ISO-8601 format)
	line := "7101f542be23e5048198e2a27c3cfda8	Title of gist item1	1 file	public	2022-09-18T18:56:10Z"

	item, err := gisty.NewGistInfo(line)
	if err != nil {
		log.Fatal(err)
	}

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
