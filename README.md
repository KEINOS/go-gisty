# go-gisty

`go-gisty` is a simple and easy-to-use **Go package for managing [GitHub Gists](https://docs.github.com/en/get-started/writing-on-github/editing-and-sharing-content-with-gists/creating-gists#about-gists)**.

## Usage

```go
go get "github.com/KIENOS/go-gisty"
```

```go
// This package requires a GitHub Personal Access Token (with gist scope)
// to be set in the GITHUB_TOKEN/GH_TOKEN environment variable.
import "github.com/KEINOS/go-gisty/gisty"
```

- [x] CRUD
  - [x] `Gisty.Create()` ..... Create a new gist with specified files to GitHub.
  - [x] `Gisty.Read()` ....... Get a content of a gist from GitHub.
  - [x] `Gisty.Update()` ..... Syncs the local changes to the gist on GitHub.
  - [x] `Gisty.Delete()` ..... Delete a specified gist from GitHub.
- [x] `Gisty.Clone()` ........ Clone a specified gist in GitHub to local.
- [x] `Gisty.List()` ......... List gists in GitHub.
- [x] `Gisty.Stargazer()` .... Get number of stars of a specified gist in GitHub.

> __Note__ : This package is a spin-off of the [`gist` subcommand](https://github.com/cli/cli/tree/trunk/pkg/cmd/gist) from the [GitHub CLI](https://docs.github.com/en/github-cli/github-cli/about-github-cli) and is intended to provide a **similar functionality as the `gh gist` command in Go applications**.
>
> Conversely, if you just **want to create a single command that perform gist operations**, it is recommended to create an alias for the `gh gist` command in your shell configuration, instead of re-inventing the wheel like I did. Also, **if you are a vim user and want to handle gist through vim**, you should consider using the [vim-gist](https://github.com/mattn/vim-gist) plugin.

```go
func Example() {
  // Create a new Gisty instance.
  obj := gisty.NewGisty()

  // The below line is equivalent to:
  //   gh gist list --public --limit 10
  args := gisty.ListArgs{
    Limit:      10,
    OnlyPublic: true,  // If both OnlyPublic and OnlySecret are true,
    OnlySecret: false, // OnlySecret is prior.
  }

  items, err := obj.List(args)
  if err != nil {
    log.Fatal(err)
  }

  // Print the fields of the first item
  firstItem := items[0]

  fmt.Println("GistID:", firstItem.GistID)
  fmt.Println("Description:", firstItem.Description)
  fmt.Println("Number of files:", firstItem.Files)
  fmt.Println("Is public gist:", firstItem.IsPublic)
  fmt.Println("Updated at:", firstItem.UpdatedAt.String())
  // Output:
  // GistID: e915aa8c01dd438e3ffd79b05f15a4ff
  // Description: Title of gist item 1
  // Number of files: 1
  // Is public gist: true
  // Updated at: 2022-04-18 03:04:38 +0000 UTC
}
```

```go
const gistIDDefault = "5b10b34f87955dfc86d310cd623a61d1"

func main() {
  // Get the gist ID from the first argument or use the default one.
  gistID := gistIDDefault

  if len(os.Args) > 1 {
    gistID = gisty.SanitizeGistID(os.Args[1])
  }

  // Instantiate a new Gisty object.
  obj := gisty.NewGisty()

  // Get the number of stars of the gist.
  count, err := obj.Stargazer(gistID)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("ID: %s, STARS: %v\n", gistID, count)
}
```

- View more examples @ pkg.go.dev
  - [function/method](https://pkg.go.dev/github.com/KEINOS/go-gisty/gisty#pkg-examples)
  - [application](https://pkg.go.dev/github.com/KEINOS/go-gisty/_examples)
