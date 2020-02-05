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
	Admin  string // contains an attribute that only admins can modify
	User   string // contains an attribute that users can modify
	Static string // contains no attributes
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
	return json.Unmarshal(def, &d)
}
