package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
)

// The Definition configures the behavior of the DNS check
// it implements the "check" interface
type Definition struct {
	Config     check.Config // generic metadata about the check
	Server     string       `optiontype:"required"`                    // The IP of the DNS server to query
	Fqdn       string       `optiontype:"required"`                    // The FQDN of the host you are looking up
	ExpectedIP string       `optiontype:"required"`                    // The expected IP of the host you are looking up
	Port       string       `optiontype:"optional" optiondefault:"53"` // The port of the DNS server
}

// Run a single instance of the check
// For now we only support A record querries
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialize empty result
	result := check.Result{Timestamp: time.Now(), Metadata: d.Config.Metadata}

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
	result.Message = "Incorrect Records Returned"
	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() check.Config {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(c check.Config) {
	d.Config = c
}
