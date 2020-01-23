package beater

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	beatcommon "github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/go-elasticsearch"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks"
	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/common"
	"gitlab.ritsec.cloud/newman/dynamicbeat/config"
	"gitlab.ritsec.cloud/newman/dynamicbeat/esclient"
)

// Dynamicbeat configuration.
type Dynamicbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	es     *elasticsearch.Client
}

// New creates an instance of dynamicbeat.
func New(b *beat.Beat, cfg *beatcommon.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	// Create the Elasticsearch client
	clientConfig := elasticsearch.Config{
		Addresses: c.CheckSource.Hosts,
		Username:  c.CheckSource.Username,
		Password:  c.CheckSource.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: !c.CheckSource.VerifyCerts, // TODO: not this
			},
		},
	}
	esClient, err := elasticsearch.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %s", err)
	}

	bt := &Dynamicbeat{
		done:   make(chan struct{}),
		config: c,
		es:     esClient,
	}
	return bt, nil
}

// Run starts dynamicbeat.
func (bt *Dynamicbeat) Run(b *beat.Beat) error {
	logp.Info("dynamicbeat is running! Hit CTRL-C to stop it.")
	var err error

	// Connect publisher client
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	// Start publisher goroutine
	pubQueue := make(chan beat.Event)
	published := make(chan uint64)
	go publishEvents(bt.client, pubQueue, published)

	// Get initial check definitions
	defs, err := esclient.UpdateCheckDefinitions(bt.es, "checks") // TODO: make check index a config
	if err != nil {
		return err
	}

	// Start running checks
	ticker := time.NewTicker(bt.config.Period)
	updateTicker := time.NewTicker(bt.config.UpdatePeriod)
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
		case <-updateTicker.C:
			// Update the check definitions
			defs, err = esclient.UpdateCheckDefinitions(bt.es, "checks") // TODO: make check index a config
			if err != nil {
				return err
			}
			logp.Info("Updated check definitions")
		case <-ticker.C:
			// Make channel for passing check definitions to and fron the checks.RunChecks goroutine
			defPass := make(chan common.CheckDefinitions)

			// Start the goroutine
			wg.Add(1)
			go checks.RunChecks(defPass, &wg, pubQueue)

			// Give it the check definitions
			defPass <- defs

			// Wait until we get the definitions back before we start the next course of checks
			defs = <-defPass
			close(defPass)
			logp.Info("Started series of checks")
		}
	}
}

func publishEvents(client beat.Client, queue <-chan beat.Event, out chan<- uint64) {
	published := uint64(0)
	for event := range queue {
		client.Publish(event)
		published++
	}
	out <- published
}

// Stop stops dynamicbeat.
func (bt *Dynamicbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func read(reader io.Reader) string {
	var buf bytes.Buffer
	buf.ReadFrom(reader)
	return buf.String()
}
