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
		checks, err := esclient.GetAllDocuments(bt.es, "checks") // TODO: factor out checks index
		if err != nil {
			return err
		}

		// Iterate over each check
		for _, check := range checks {
			// Get any template variables for the check
			attribs := make(map[string]string)
			id := check["id"].String()
			for _, perm := range []string{"admin", "user"} {
				// Generate attribute index name
				idx := strings.Join([]string{"attrib_", perm, "_", id}, "")

				// Get attribute document
				docs, err := esclient.GetAllDocuments(bt.es, idx)
				if err != nil {
					return err
				}
				attrib := docs[0]

				// Read attributes from document
				for k, v := range attrib {
					if _, pres := attribs[k]; !pres {
						attribs[k] = v.String()
					}
				}
			}

			// Template out any mapping values
			type NoopAttributes struct {
				AdminName, TeamName string
			}
			templAttribs := NoopAttributes{
				AdminName: attribs["AdminName"],
				TeamName:  attribs["TeamName"],
			}
			def := make(map[string]string)
			for k, v := range check["definition"].Map() {
				templ := template.Must(template.New(k).Parse(v.String()))
				var buf bytes.Buffer
				if err := templ.Execute(&buf, templAttribs); err != nil {
					return fmt.Errorf("Error parsing template for key %s: %s", k, err)
				}
				def[k] = buf.String()
			}

			// Send message
			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type": b.Info.Name,
					"message": strings.Join([]string{
						def["static"],
						def["admin"],
						def["team"],
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
