package run

import (
	"context"
	"sync"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/models"
	"go.uber.org/zap"
)

// Round : Run a course of checks based on the currently-loaded configuration.
func Round(defs []models.CheckConfig, results chan<- check.Result, started chan<- bool) {
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
		names[d.CheckId] = false
		wg.Add(1)

		def := d
		go func() {
			defer wg.Done()

			checkStart := time.Now()
			result := Check(ctx, def)
			zap.S().Debugf("[%s] Finished after %.2f seconds", result.CheckId, time.Since(checkStart).Seconds())
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
			if len(names) == 0 {
				break
			} else {
				time.Sleep(30 * time.Second)
				if len(names) > 0 {
					zap.S().Warnf("Checks still running after %.2f seconds: %+v", time.Since(start).Seconds(), names)
				}
			}
		}
		zap.S().Infof("All checks started %.2f seconds ago have finished", time.Since(start).Seconds())
		close(finished)
	}()
	for result := range finished {
		// Record that the check has finished
		delete(names, result.CheckId)

		// Publish the event to the publisher queue
		results <- result
	}
}
