package icmp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sparrc/go-ping"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the ICMP check
// it implements the "Check" interface
type Definition struct {
	Config    schema.CheckConfig // generic metadata about the check
	Host      string             `optiontype:"required"`                       // IP or hostname of the host to run the ICMP check against
	Count     int                `optiontype:"optional" optiondefault:"1"`     // The number of ICMP requests to send per check
	PassCount string             `optiontype:"optional" optionadefault:"true"` // Pass check based on received pings matching Count; if false, will use percent packet loss
	Percent   int                `optiontype:"optional" optiondefault:"100"`   // The percent threshold of dropped packets to fail a check
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Create pinger
	pinger, err := ping.NewPinger(d.Host)
	if err != nil {
		result.Message = fmt.Sprintf("Error creating pinger: %s", err)
		return result
	}

	// Send ping
	pinger.Count = 3
	// TODO: change this to be relative to the parent context's timeout
	pinger.Timeout = 25 * time.Second
	pinger.Run()

	// Convert PassCount to bool
	passCount, err := strconv.ParseBool(d.PassCount)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to parse PassCount boolean from struct def : %s", err)
		return result
	}

	stats := pinger.Statistics()

	// Check packet loss instead of count
	if !passCount {
		if stats.PacketLoss >= float64(d.Percent) {
			result.Message = fmt.Sprint("Not all pings made it back!")
			return result
		}

		// If we make it here the check passes by percentage
		result.Passed = true
		return result
	}

	// Check for failure of ICMP
	if stats.PacketsRecv != d.Count {
		result.Message = fmt.Sprint("Not all pings made it back!")
		return result
	}

	// If we make it here the check passes
	result.Passed = true
	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
