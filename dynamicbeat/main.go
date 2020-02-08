package main

import (
	"os"

	"github.com/newman/scorestack/dynamicbeat/cmd"

	_ "github.com/newman/scorestack/dynamicbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
