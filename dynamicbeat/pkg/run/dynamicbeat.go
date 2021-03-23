package run

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/scorestack/scorestack/dynamicbeat/checks"
	"github.com/scorestack/scorestack/dynamicbeat/checks/schema"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
)

const index = "checkdef"

// Run starts dynamicbeat.
func Run() error {
	// logp.Info("dynamicbeat is running! Hit CTRL-C to stop it.")
	c := config.Get()

	// Create the Elasticsearch client
	clientConfig := elasticsearch.Config{
		Addresses: []string{c.Elasticsearch},
		Username:  c.Username,
		Password:  c.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			DialContext:         (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !c.VerifyCerts,
			},
		},
	}
	es, err := elasticsearch.NewClient(clientConfig)
	if err != nil {
		return fmt.Errorf("Error creating client: %s", err)
	}

	// Connect publisher client
	/*
		bt.client, err := b.Publisher.Connect()
		if err != nil {
			return err
		}
	*/

	// Get initial check definitions
	var defs []schema.CheckConfig
	doubleBreak := false
	// TODO: Find a better way for looping until we can hit Elasticsearch
	// logp.Info("Getting initial check definitions...")
	for {
		select {
		// Case for catching Ctrl+C and gracfully exiting
		case <-bt.done:
			return nil
		default:
			// Continue looping and sleeping till we can hit Elasticsearch
			defs, err = esclient.UpdateCheckDefs(es, index)
			if err != nil {
				// logp.Info("Failed to reach Elasticsearch. Waiting 5 seconds to try again...")
				// logp.Debug("dynamicbeat", "Connection error was: %s", err)
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
	pubQueue := make(chan beat.Event)
	published := make(chan uint64)
	go publishEvents(bt.client, pubQueue, published)

	// Start running checks
	ticker := time.NewTicker(c.RoundTime)

	var wg sync.WaitGroup
	for {
		select {
		case <-bt.done:
			// Wait for all checks.RunChecks goroutines to exit
			wg.Wait()

			// Close the event publishing queue so the publishEvents goroutine will exit
			close(pubQueue)

			// Wait for all events to be published
			<-published
			close(published)
			return nil
		case <-ticker.C:
			// logp.Info("Number of go-routines: %d", runtime.NumGoroutine())
			// logp.Info("Starting a series of %d checks", len(defs))

			// Make channel for passing check definitions to and fron the checks.RunChecks goroutine
			defPass := make(chan []schema.CheckConfig)

			// Start the goroutine
			wg.Add(1)
			go func() {
				defer wg.Done()
				checks.RunChecks(defPass, pubQueue)
			}()

			// Give it the check definitions
			defPass <- defs

			// Wait until we get the definitions back before we start the next course of checks
			defs = <-defPass
			close(defPass)
			// logp.Info("Started series of checks")

			// Update the check definitions for the next round
			defs, err = esclient.UpdateCheckDefs(es, index)
			if err != nil {
				// logp.Info("Failed to update check definitions : %s", err)
			}
		}
	}
}

func publishEvents(es elasticsearch.Client, queue <-chan beat.Event, out chan<- uint64) {
	published := uint64(0)
	for event := range queue {
		client.Publish(event)
		published++
	}
	out <- published
}
