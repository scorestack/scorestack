package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/s-newman/scorestack/dynamicbeat/cmd"

	_ "github.com/s-newman/scorestack/dynamicbeat/include"
)

func run() error {
	// defer profile.Start(profile.ThreadcreationProfile).Stop()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	return cmd.RootCmd.Execute()
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
