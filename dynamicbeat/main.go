package main

import (
	"os"

	"github.com/pkg/profile"

	"github.com/s-newman/scorestack/dynamicbeat/cmd"

	_ "github.com/s-newman/scorestack/dynamicbeat/include"
)

func run() error {
	defer profile.Start(profile.GoroutineProfile).Stop()
	return cmd.RootCmd.Execute()
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
