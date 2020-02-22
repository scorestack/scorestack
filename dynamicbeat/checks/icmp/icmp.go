package icmp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sparrc/go-ping"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the ICMP check
// it implements the "Check" interface
type Definition struct {
	ID          string  // unique identifier for this check
	Name        string  // a human-readable title for this check
	Group       string  // the group this check is part of
	ScoreWeight float64 // the weight that this check has relative to others
	Host        string  // (required) IP or hostname of the host to run the ICMP check against
	Count       int     // (opitonal, default=1) The number of ICMP requests to send per check
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context, wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "icmp",
	}

	// Make channels for completing a check or not
	done := make(chan bool)
	failed := make(chan bool)

	go func() {
		// Create pinger
		pinger, err := ping.NewPinger(d.Host)
		if err != nil {
			result.Message = fmt.Sprintf("Error creating pinger: %s", err)
			failed <- true
			return
		}

		// Send ping
		pinger.Count = d.Count
		// pinger.Timeout = 5 * time.Second
		pinger.Run()

		stats := pinger.Statistics()

		// Check for failure of ICMP
		if stats.PacketsRecv != d.Count {
			result.Message = fmt.Sprintf("FAILED: Not all pings made it back! Received %d out of %d", stats.PacketsRecv, stats.PacketsSent)
			failed <- true
			return
		}

		// If we make it here the check passes
		done <- true
	}()

	// Watch channels and context for timeout
	for {
		select {
		case <-done:
			result.Passed = true
			out <- result
			return
		case <-failed:
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

	// Explicitly set default values
	d.Count = 1

	// Unpack json definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set generic attributes
	d.ID = id
	d.Name = name
	d.Group = group
	d.ScoreWeight = scoreWeight

	// Make sure required fields are defined
	missingFields := make([]string, 0)
	if d.Host == "" {
		missingFields = append(missingFields, "Host")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "icmp",
			Field: missingFields[0],
		}
	}
	return nil
}
