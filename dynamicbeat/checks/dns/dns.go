package dns

import (
	"encoding/json"
	"sync"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the DNS check
// it implements the "check" interface
type Definition struct {
	ID         string // a unique identifier for this check
	Name       string // a human-readable title for the check
	Group      string // the group this check is part of
	Server     string // (required) The IP of the DNS server to query
	Fqdn       string // (required) The FQDN of the host you are looking up
	ExpectedIP string // (required) The expected IP of the host you are looking up
	Port       string // (optional, default=53) The port of the DNS server
}

// Run a single instance of the check
// For now we only support A record querries
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {

}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, def []byte) error {

	// Explicitly set default values
	d.Port = "53"

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
	missingFields := make([]string, 0)
	if d.Server == "" {
		missingFields = append(missingFields, "Server")
	}

	if d.Fqdn == "" {
		missingFields = append(missingFields, "Fqdn")
	}

	if d.ExpectedIP == "" {
		missingFields = append(missingFields, "ExpectedIP")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "dns",
			Field: missingFields[0],
		}
	}
	return nil
}
