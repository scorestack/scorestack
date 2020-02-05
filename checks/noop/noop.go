package noop

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of a Noop check.
type Definition struct {
	ID     string // a unique identifier for this check
	Name   string // a human-readable title for this check
	Admin  string // (required) contains an attribute that only admins can modify
	User   string // (required) contains an attribute that users can modify
	Static string // (required) contains no attributes
}

// Run a single instance of the check.
func (d Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	result := schema.CheckResult{
		Timestamp: time.Now(),
		ID:        d.ID,
		Name:      d.Name,
		CheckType: "noop",
		Passed:    true,
		Message:   strings.Join([]string{d.Admin, d.User, d.Static}, "; "),
		Details: map[string]string{
			"Admin":  d.Admin,
			"User":   d.User,
			"Static": d.Static,
		},
	}

	out <- result
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d Definition) Init(id string, name string, def []byte) error {
	// Set ID and name attributes
	d.ID = id
	d.Name = name

	// Unpack definition json
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Make sure required fields are defined
	missingFields := make([]string, 0)
	if d.Admin == "" {
		missingFields = append(missingFields, "Admin")
	}

	if d.User == "" {
		missingFields = append(missingFields, "User")
	}

	if d.Static == "" {
		missingFields = append(missingFields, "Static")
	}

	// Error on only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "noop",
			Field: missingFields[0],
		}
	}
	return nil
}
