package run

import (
	"context"
	"sync"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"go.uber.org/zap"
)

// Round : Run a course of checks based on the currently-loaded configuration.
func Round(defs []check.Config, results chan<- check.Result, started chan<- bool) {
	start := time.Now()

	// Make an event queue separate from the publisher queue so we can track
	// which checks are still running
	finished := make(chan check.Result, len(defs))

	// Iterate over each check
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	names := make(map[string]bool)
	var wg sync.WaitGroup
	for _, d := range defs {
		// Start check goroutine
		names[d.ID] = false
		wg.Add(1)

		def := d
		go func() {
			defer wg.Done()

			checkStart := time.Now()
			result := Check(ctx, def)
			zap.S().Infof("[%s] Finished after %.2f seconds", result.ID, time.Since(checkStart).Seconds())
			finished <- result
		}()
	}

	// Signal that all checks have started
	started <- true

	// Wait for checks to finish
	defer wg.Wait()
	// zap.S().Infof("Checks started at %s have finished in %.2f seconds", start.Format("15:04:05.000"), time.Since(start).Seconds())
	go func() {
		for {
			if names == nil {
				break
			} else if len(names) == 0 {
				break
			} else {
				time.Sleep(30 * time.Second)
				zap.S().Infof("Checks still running after %.2f seconds: %+v", time.Since(start).Seconds(), names)
			}
		}
		zap.S().Infof("All checks started %.2f seconds ago have finished", time.Since(start).Seconds())
		close(finished)
	}()
	for result := range finished {
		// Record that the check has finished
		delete(names, result.ID)

		// Publish the event to the publisher queue
		results <- result
	}
}
