package main

import (
	"log"

	"github.com/scorestack/scorestack/dynamicbeat/cmd"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	cmd.Execute()
}
