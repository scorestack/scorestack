package noop

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of a Noop check.
type Definition struct {
	ID          string  // a unique identifier for this check
	Name        string  // a human-readable title for this check
	Group       string  // the group this check is part of
	ScoreWeight float64 // the weight that this check has relative to others
	Dynamic     string  // (required) contains attributes that can be modified by admins or users
	Static      string  // (required) contains no attributes
}

// Run a single instance of the check.
func (d *Definition) Run(ctx context.Context, wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "noop",
		Passed:      true,
		Message:     strings.Join([]string{d.Dynamic, d.Static}, "; "),
		Details: map[string]string{
			"Dynamic": d.Dynamic,
			"Static":  d.Static,
		},
	}

	// make channel for completing the check or not
	done := make(chan bool)

	go func() {
		done <- true
		return
	}()

	// Watch channels and context for timeout
	for {
		select {
		case <-done:
			close(done)
			result.Passed = true
			out <- result
			return
		case <-ctx.Done():
			result.Message = fmt.Sprintf("Timeout via context : %s", ctx.Err())
			out <- result
			return
		}
	}
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, scoreWeight float64, def []byte) error {
	// Unpack definition json
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set ID and name attributes
	d.ID = id
	d.Name = name
	d.Group = group
	d.ScoreWeight = scoreWeight

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
