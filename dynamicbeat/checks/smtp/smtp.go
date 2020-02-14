package smtp

import (
	"encoding/json"
	"sync"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SMTP check
// it implements the "check" interface
type Definition struct {
	ID       string // a unique identifier for this check
	Name     string // a human-readable title for the check
	Group    string // the group this check is part of
	Host     string // (required) IP or hostname of the smtp server
	Username string // (required) Username for the smtp server
	Password string // (required) Password for the smtp server
	Sender   string // (required) Who is sending the email
	Reciever string // (required) Who is receiving the email
	Body     string // (optional, default="Hello from scoring engine") Body of the email
	Port     string // (optional, default="25") Port of the smtp server
}

// Run a single instance of the check
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {

}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, def []byte) error {

	// Explicitly set defaults
	d.Body = "Hello from scoring engine"
	d.Port = "25"

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
	if d.Host == "" {
		missingFields = append(missingFields, "Host")
	}

	if d.Username == "" {
		missingFields = append(missingFields, "Username")
	}

	if d.Password == "" {
		missingFields = append(missingFields, "Password")
	}

	if d.Sender == "" {
		missingFields = append(missingFields, "Sender")
	}

	if d.Reciever == "" {
		missingFields = append(missingFields, "Reciever")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "smtp",
			Field: missingFields[0],
		}
	}
	return nil
}
