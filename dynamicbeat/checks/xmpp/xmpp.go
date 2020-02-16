package xmpp

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
	"gosrc.io/xmpp"
)

// The Definition configures the behavior of the XMPP check
// it implements the "check" interface
type Definition struct {
	ID        string // a unique identifier for this check
	Name      string // a human-readable title for the check
	Group     string // the group this check is part of
	Host      string // (required) IP or hostname of the xmpp server
	Username  string // (required) Username to use for the xmpp server
	Password  string // (required) Password for the user
	Encrypted bool   // (optional, default=true) TLS support or not
	Port      string // (optional, default=5222) Port for the xmpp server
}

// Run a single instance of the check
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	// Set up result
	result := schema.CheckResult{
		Timestamp: time.Now(),
		ID:        d.ID,
		Name:      d.Name,
		Group:     d.Group,
		CheckType: "xmpp",
	}

	// Create xmpp config
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: fmt.Sprintf("%s:%s", d.Host, d.Port),
		},
	}

}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, def []byte) error {

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
