package main

import (
	"os"

	"github.com/s-newman/scorestack/dynamicbeat/cmd"

	_ "github.com/s-newman/scorestack/dynamicbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
