package main

import (
	"os"

	"gitlab.ritsec.cloud/newman/dynamicbeat/cmd"

	_ "gitlab.ritsec.cloud/newman/dynamicbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
