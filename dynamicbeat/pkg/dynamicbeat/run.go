package dynamicbeat

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checksource"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/run"
	"go.uber.org/zap"
)

const CHECKDEF_INDEX = "checkdef"

// Run starts dynamicbeat.
func Run(rounds uint) error {
	zap.S().Infof("dynamicbeat is running! Hit CTRL-C to stop it.")
	c := config.Get()

	pub, err := esclient.New(c.Elasticsearch, c.Username, c.Password, c.VerifyCerts)
	if err != nil {
		return err
	}

	es, err := checksource.NewElasticsearch(c.Elasticsearch, c.Username, c.Password, c.VerifyCerts, CHECKDEF_INDEX)
	if err != nil {
		return err
	}

	// Connect publisher client
	/*
		bt.client, err := b.Publisher.Connect()
		if err != nil {
			return err
		}
	*/

	// Set up CTRL+C handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Get initial check definitions
	var defs []check.Config
	doubleBreak := false
	// TODO: Find a better way for looping until we can hit Elasticsearch
	zap.S().Infof("Getting initial check definitions...")
	for {
		select {
		// Case for catching Ctrl+C and gracfully exiting
		case <-quit:
			return nil
		default:
			// Continue looping and sleeping till we can hit Elasticsearch
			defs, err = es.LoadAll()
			if err != nil {
				zap.S().Infof("Failed to reach Elasticsearch. Waiting 5 seconds to try again...")
				zap.S().Debugf("dynamicbeat", "Connection error was: %s", err)
				time.Sleep(5 * time.Second)
			} else {
				doubleBreak = true
				break
			}
		}
		// TODO: Find a better way of breaking the for loop if we break from switch
		// Needed to break out of the for loop
		if doubleBreak {
			break
		}
	}

	// Start publisher goroutine
	results := make(chan check.Result)
	published := make(chan uint64)
	go publishEvents(pub, results, published)

	// Start round timer
	ticker := time.NewTicker(c.RoundTime)

	// Start first round
	var wg sync.WaitGroup
	runRound(defs, results, &wg)
	roundsRun := uint(1)

	if rounds == 1 {
		return waitAndExit(&wg, results, published)
	}

	for {
		select {
		case <-quit:
			return waitAndExit(&wg, results, published)
		case <-ticker.C:
			zap.S().Infof("Number of goroutines: %d", runtime.NumGoroutine())
			runRound(defs, results, &wg)
			roundsRun += 1

			if rounds > 0 && roundsRun >= rounds {
				return waitAndExit(&wg, results, published)
			}

			// Update the check definitions for the next round
			defs, err = es.LoadAll()
			if err != nil {
				zap.S().Warnf("Failed to update check definitions : %s", err)
			}
		}
	}
}

func runRound(defs []check.Config, results chan<- check.Result, wg *sync.WaitGroup) {
	zap.S().Infof("Starting a series of %d checks", len(defs))

	// Start the goroutine
	started := make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		run.Round(defs, results, started)
	}()

	// Wait until all the checks have been started before we refresh
	// the checks from Elasticsearch to make sure that we don't
	// overwrite the check definitions while they're in use.
	<-started
	zap.S().Infof("Started series of checks")
}

func waitAndExit(wg *sync.WaitGroup, results chan check.Result, published chan uint64) error {
	// Wait for all checks.RunChecks goroutines to exit
	wg.Wait()

	// Close the event publishing queue so the publishEvents goroutine will exit
	close(results)

	// Wait for all events to be published
	<-published
	close(published)
	return nil
}

func publishEvents(es *esclient.Client, results <-chan check.Result, out chan<- uint64) {
	published := uint64(0)
	for result := range results {
		err := es.AddResult(result)
		if err != nil {
			zap.S().Error(err)
			zap.S().Errorf("check that failed to index: %+v", result)
		} else {
			published++
		}
	}
	out <- published
}
