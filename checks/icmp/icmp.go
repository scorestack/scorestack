package icmp

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sparrc/go-ping"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the ICMP check
// it implements the "Check" interface
type Definition struct {
	ID    string // unique identifier for this check
	Name  string // a human-readable title for this check
	IP    string // (required) IP of the host to run the ICMP check against
	Count int    // (opitonal, default=1) The number of ICMP requests to send per check
}

// Run a single instance of the check
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	// Set up result
	result := schema.CheckResult{
		Timestamp: time.Now(),
		ID:        d.ID,
		Name:      d.Name,
		CheckType: "icmp",
	}

	// Create pinger
	pinger, err := ping.NewPinger(d.IP)
	if err != nil {
		result.Message = fmt.Sprintf("Error creating pinger: %s", err)
		out <- result
	}

	// Send ping
	pinger.Count = d.Count
	pinger.Timeout = 5 * time.Second
	pinger.Run()

	stats := pinger.Statistics()

	// Check for failure of ICMP
	if stats.PacketsRecv != d.Count {
		result.Message = fmt.Sprintf("FAILED: Not all pings made it back! Received %d out of %d", stats.PacketsRecv, stats.PacketsSent)
		out <- result
		return
	}

	// If we make it here the check passes
	result.Message = fmt.Sprintf("SUCCESS")
	out <- result
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, def []byte) error {

	// Unpack json definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Check for optional value being set
	if d.Count <= 0 {
		d.Count = 1
	}

	// Set ID and name attributes
	d.ID = id
	d.Name = name

	// Make sure required fields are defined
	missingFields := make([]string, 0)
	if d.IP == "" {
		missingFields = append(missingFields, "IP")
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
