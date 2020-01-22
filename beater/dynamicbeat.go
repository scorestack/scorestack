package beater

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"text/template"
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
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		// Get list of checks
		resp, err := bt.es.Search(
			bt.es.Search.WithIndex("checks"), // TODO: factor out checks index
			bt.es.Search.WithStoredFields("_id"),
		)
		if err != nil {
			return fmt.Errorf("Error getting check IDs: %s", err)
		}
		defer resp.Body.Close()
		respBody := read(resp.Body)
		checkIDs := gjson.Get(respBody, "hits.hits.#._id").Array()

		// Iterate over each check
		for _, checkID := range checkIDs {
			// Get the check document
			resp, err := bt.es.Get("checks", checkID.String())
			if err != nil {
				return fmt.Errorf("Error getting check %s: %s", checkID.String(), err)
			}
			defer resp.Body.Close()
			respBody = read(resp.Body)

			// Get any template variables for the check
			var attributes map[string]string
			id := gjson.Get(respBody, "id").String()
			for _, indexType := range []string{"admin", "user"} {
				// Generate attribute index name
				indexName := strings.Join([]string{"attrib_", indexType, "_", id}, "")

				// Find index's ID
				resp, err := bt.es.Search(
					bt.es.Search.WithIndex(indexName),
					bt.es.Search.WithStoredFields("_id"),
				)
				if err != nil {
					return fmt.Errorf("Error getting attribute index %s: %s", indexName, err)
				}
				defer resp.Body.Close()
				id := gjson.Get(read(resp.Body), "hits.hits.0._id").String()

				// Get attribute document
				resp, err = bt.es.Get(indexName, id)
				if err != nil {
					return fmt.Errorf("Error reading attribute document %s: %s", id, err)
				}
				defer resp.Body.Close()

				// Read attributes from document
				for key, val := range gjson.Get(read(resp.Body), "*").Map() {
					if _, present := attributes[key]; !present {
						attributes[key] = val.String()
					}
				}
			}

			// Template out any mapping values
			type NoopAttributes struct {
				AdminName, TeamName string
			}
			templateAttributes := NoopAttributes{
				AdminName: attributes["AdminName"],
				TeamName:  attributes["TeamName"],
			}
			var templatedDefinition map[string]string
			for key, val := range gjson.Get(respBody, "definition").Map() {
				valTemplate := template.Must(template.New(key).Parse(val.String()))
				var buf bytes.Buffer
				if err := valTemplate.Execute(&buf, templateAttributes); err == nil {
					return fmt.Errorf("Error parsing template for key %s: %s", key, err)
				}
				templatedDefinition[key] = buf.String()
			}

			// Send message
			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type": b.Info.Name,
					"message": strings.Join([]string{
						templatedDefinition["static"],
						templatedDefinition["admin"],
						templatedDefinition["team"],
					}, " - "),
				},
			}
			bt.client.Publish(event)
			logp.Info("Event sent")
		}
	}
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
