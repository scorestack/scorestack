package main

import (
	"os"

	"github.com/scorestack/scorestack/dynamicbeat/cmd"

	_ "github.com/scorestack/scorestack/dynamicbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
