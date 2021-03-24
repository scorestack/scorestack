package run

import (
	"context"
	"sync"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/event"
	"go.uber.org/zap"
)

// Round : Run a course of checks based on the currently-loaded configuration.
func Round(defs []check.Config, results chan<- event.Event, started chan<- bool) {
	start := time.Now()

	// Make an event queue separate from the publisher queue so we can track
	// which checks are still running
	eventQueue := make(chan event.Event, len(defs))

	// Iterate over each check
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	names := make(map[string]bool)
	var wg sync.WaitGroup
	for _, def := range defs {
		names[def.Meta.ID] = false
		chk, err := unpackDef(def)
		if err != nil {
			// Something was wrong with templating the check. Return a failed event with the error.
			errorDetail := make(map[string]string)
			errorDetail["error_message"] = err.Error()
			eventQueue <- event.Event{
				Timestamp:   time.Now(),
				Id:          chk.GetConfig().Meta.ID,
				Name:        chk.GetConfig().Meta.Name,
				CheckType:   chk.GetConfig().Meta.Type,
				Group:       chk.GetConfig().Meta.Group,
				ScoreWeight: chk.GetConfig().Meta.ScoreWeight,
				Passed:      false,
				Message:     "Encountered an error when unpacking check definition.",
				Details:     errorDetail,
			}
		}

		// Start check goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()

			checkStart := time.Now()
			checkName := chk.GetConfig().Meta.Name
			eventQueue <- runCheck(ctx, chk)
			zap.S().Infof("[%s] Finished after %.2f seconds", checkName, time.Since(checkStart).Seconds())
		}()
	}

	// Signal that all checks have started
	started <- true

	// Wait for checks to finish
	defer wg.Wait()
	zap.S().Infof("Checks started at %s have finished in %.2f seconds", start.Format("15:04:05.000"), time.Since(start).Seconds())
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
		close(eventQueue)
	}()
	for evt := range eventQueue {
		// Record that the check has finished
		delete(names, evt.Id)

		// Publish the event to the publisher queue
		results <- evt
	}
}
