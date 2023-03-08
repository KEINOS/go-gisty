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
