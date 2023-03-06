package main

import (
	"fmt"
	"log"

	"github.com/KENOS/go-gisty/gisty"
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
