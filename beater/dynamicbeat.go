package beater

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/go-elasticsearch"
	"github.com/tidwall/gjson"

	"gitlab.ritsec.cloud/newman/dynamicbeat/config"
)

// Dynamicbeat configuration.
type Dynamicbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	es     *elasticsearch.Client
}

// New creates an instance of dynamicbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
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
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		resp, err := bt.es.Search(
			bt.es.Search.WithIndex("checks"),
			bt.es.Search.WithStoredFields("_id"),
		)
		if err != nil {
			return fmt.Errorf("Error getting check IDs: %s", err)
		}
		defer resp.Body.Close()
		var byteReader bytes.Buffer
		byteReader.ReadFrom(resp.Body)
		respBody := byteReader.String()
		hits := gjson.Get(respBody, "hits.total.value").Num

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": counter,
				"hits":    hits,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
		counter++
	}
}

// Stop stops dynamicbeat.
func (bt *Dynamicbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
