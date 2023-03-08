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
