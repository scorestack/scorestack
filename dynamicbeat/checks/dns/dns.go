package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the DNS check
// it implements the "check" interface
type Definition struct {
	Config     schema.CheckConfig // generic metadata about the check
	Server     string             // (required) The IP of the DNS server to query
	Fqdn       string             // (required) The FQDN of the host you are looking up
	ExpectedIP string             // (required) The expected IP of the host you are looking up
	Port       string             // (optional, default=53) The port of the DNS server
}

// Run a single instance of the check
// For now we only support A record querries
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Setup for dns query
	var msg dns.Msg
	fqdn := dns.Fqdn(d.Fqdn)
	msg.SetQuestion(fqdn, dns.TypeA)

	// Make it obey timeout via deadline
	// TODO: change this to be relative to the parent context's timeout
	deadctx, cancel := context.WithDeadline(ctx, time.Now().Add(20*time.Second))
	defer cancel()

	// Send the query
	in, err := dns.ExchangeContext(deadctx, &msg, fmt.Sprintf("%s:%s", d.Server, d.Port))
	if err != nil {
		result.Message = fmt.Sprintf("Problem sending query to %s : %s", d.Server, err)
		return result
	}

	// Check if we got any records
	if len(in.Answer) < 1 {
		result.Message = fmt.Sprintf("No records received from %s", d.Server)
		return result
	}

	// Loop through results and check for correct match
	for _, answer := range in.Answer {
		// Check if an answer is an A record and it matches the expected IP
		if a, ok := answer.(*dns.A); ok && (a.A).String() == d.ExpectedIP {
			// If we reach here the check succeeds
			result.Passed = true
			return result
		}
	}

	// If we reach here no records matched expected IP and check fails
	result.Message = fmt.Sprintf("Incorrect Records Returned")
	return result
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(config schema.CheckConfig, def []byte) error {

	// Explicitly set default values
	d.Port = "53"

	// Unpack JSON definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set generic values
	d.Config = config

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
			ID:    d.Config.ID,
			Type:  "dns",
			Field: missingFields[0],
		}
	}
	return nil
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}
