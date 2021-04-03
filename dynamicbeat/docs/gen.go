package main

import (
	"log"

	"github.com/scorestack/scorestack/dynamicbeat/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	dynamicbeat := cmd.NewRootCommand()
	err := doc.GenMarkdownTree(dynamicbeat, "../docs/dynamicbeat/reference")
	if err != nil {
		log.Fatal(err)
	}
}
