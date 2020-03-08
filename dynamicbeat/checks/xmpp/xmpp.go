package xmpp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"gosrc.io/xmpp/stanza"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
	"gosrc.io/xmpp"
)

// The Definition configures the behavior of the XMPP check
// it implements the "check" interface
type Definition struct {
	Config    schema.CheckConfig // generic metadata about the check
	Host      string             // (required) IP or hostname of the xmpp server
	Username  string             // (required) Username to use for the xmpp server
	Password  string             // (required) Password for the user
	Encrypted bool               // (optional, default=true) TLS support or not
	Port      string             // (optional, default=5222) Port for the xmpp server
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.Config.ID,
		Name:        d.Config.Name,
		Group:       d.Config.Group,
		ScoreWeight: d.Config.ScoreWeight,
		CheckType:   "xmpp",
	}

	// Create xmpp config
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address:   fmt.Sprintf("%s:%s", d.Host, d.Port),
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
			// ConnectTimeout: 20,
		},
		Jid:        fmt.Sprintf("%s@%s", d.Username, d.Host),
		Credential: xmpp.Password(d.Password),
		Insecure:   d.Encrypted,
		// ConnectTimeout: 20,
	}

	// KILL THIS CHECK!
	passed := make(chan bool)
	failed := make(chan bool)

	go func() {
		// Create a client
		client, err := xmpp.NewClient(config, xmpp.NewRouter(), errorHandler)
		if err != nil {
			result.Message = fmt.Sprintf("Creating a xmpp client failed : %s", err)
			failed <- true
			return
			// return result
		}

		// Create IQ xmpp message
		iq, err := stanza.NewIQ(stanza.Attrs{
			Type: stanza.IQTypeGet,
			From: d.Host,
			To:   "localhost",
			Id:   "ScoreStack-check",
		})
		if err != nil {
			result.Message = fmt.Sprintf("Creating IQ message failed : %s", err)
			failed <- true
			return
			// return result
		}

		// Set Disco as the payload of IQ
		disco := iq.DiscoInfo()
		iq.Payload = disco

		// Connect the client
		err = client.Connect()
		if err != nil {
			result.Message = fmt.Sprintf("Connecting to %s failed : %s", d.Host, err)
			failed <- true
			return
			// return result
		}
		defer func() {
			if closeErr := client.Disconnect(); closeErr != nil {
				// logp.Warn("failed to close xmpp connection: %s", closeErr.Error())
			}
		}()

		// Send the IQ message
		err = client.Send(iq)
		if err != nil {
			result.Message = fmt.Sprintf("Sending IQ message to %s failed %s", d.Host, err)
			failed <- true
			return
			// return result
		}

		// If we make it here the check should pass
		result.Passed = true
		passed <- true
		return
		// return result
	}()

	for {
		select {
		case <-ctx.Done():
			result.Message = fmt.Sprintf("Timeout limit reached: %s", ctx.Err())
			return result
		case <-failed:
			return result
		case <-passed:
			return result
		}
	}
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(config schema.CheckConfig, def []byte) error {

	// Set explicit values
	d.Port = "5222"
	d.Encrypted = true

	// Unpack JSON definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set generic values
	d.Config = config

	// Check for missing fields
	missingfields := make([]string, 0)
	if d.Host == "" {
		missingfields = append(missingfields, "Host")
	}

	if d.Username == "" {
		missingfields = append(missingfields, "Username")
	}

	if d.Password == "" {
		missingfields = append(missingfields, "Password")
	}

	// Error only the first missing field, if there are any
	if len(missingfields) > 0 {
		return schema.ValidationError{
			ID:    d.Config.ID,
			Type:  "xmpp",
			Field: missingfields[0],
		}
	}
	return nil
}

// Without this function, the xmpp "client" calls will seg fault
func errorHandler(err error) {
	return
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}
