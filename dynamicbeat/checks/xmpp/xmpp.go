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
	ID          string  // a unique identifier for this check
	Name        string  // a human-readable title for the check
	Group       string  // the group this check is part of
	ScoreWeight float64 // the weight that this check has relative to others
	Host        string  // (required) IP or hostname of the xmpp server
	Username    string  // (required) Username to use for the xmpp server
	Password    string  // (required) Password for the user
	Encrypted   bool    // (optional, default=true) TLS support or not
	Port        string  // (optional, default=5222) Port for the xmpp server
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "xmpp",
	}

	// Create xmpp config
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address:        fmt.Sprintf("%s:%s", d.Host, d.Port),
			TLSConfig:      &tls.Config{InsecureSkipVerify: true},
			ConnectTimeout: 20,
		},
		Jid:            fmt.Sprintf("%s@%s", d.Username, d.Host),
		Credential:     xmpp.Password(d.Password),
		Insecure:       d.Encrypted,
		ConnectTimeout: 20,
	}

	// Create a client
	client, err := xmpp.NewClient(config, xmpp.NewRouter(), errorHandler)
	if err != nil {
		result.Message = fmt.Sprintf("Creating a xmpp client failed : %s", err)
		return result
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
		return result
	}

	// Set Disco as the payload of IQ
	disco := iq.DiscoInfo()
	iq.Payload = disco

	// Connect the client
	err = client.Connect()
	if err != nil {
		result.Message = fmt.Sprintf("Connecting to %s failed : %s", d.Host, err)
		return result
	}
	defer client.Disconnect()

	// Send the IQ message
	err = client.Send(iq)
	if err != nil {
		result.Message = fmt.Sprintf("Sending IQ message to %s failed %s", d.Host, err)
		return result
	}

	// If we make it here the check should pass
	result.Passed = true
	return result
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, scoreWeight float64, def []byte) error {

	// Set explicit values
	d.Port = "5222"
	d.Encrypted = true

	// Unpack JSON definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set generic values
	d.ID = id
	d.Name = name
	d.Group = group
	d.ScoreWeight = scoreWeight

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
			ID:    d.ID,
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
