package noop

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of a Noop check.
type Definition struct {
	ID      string // a unique identifier for this check
	Name    string // a human-readable title for this check
	Dynamic string // (required) contains attributes that can be modified by admins or users
	Static  string // (required) contains no attributes
}

// Run a single instance of the check.
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	result := schema.CheckResult{
		Timestamp: time.Now(),
		ID:        d.ID,
		Name:      d.Name,
		CheckType: "noop",
		Passed:    true,
		Message:   strings.Join([]string{d.Dynamic, d.Static}, "; "),
		Details: map[string]string{
			"Dynamic": d.Dynamic,
			"Static":  d.Static,
		},
	}

	out <- result
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, def []byte) error {
	// Unpack definition json
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set ID and name attributes
	d.ID = id
	d.Name = name

	// Make sure required fields are defined
	missingFields := make([]string, 0)
	if d.Dynamic == "" {
		missingFields = append(missingFields, "Dynamic")
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