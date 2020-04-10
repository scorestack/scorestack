package icmp

import (
	"context"
	"fmt"
	"time"

	"github.com/sparrc/go-ping"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the ICMP check
// it implements the "Check" interface
type Definition struct {
	Config schema.CheckConfig // generic metadata about the check
	Host   string             `optiontype:"required"`                   // IP or hostname of the host to run the ICMP check against
	Count  int                `optiontype:"optional" optiondefault:"1"` // The number of ICMP requests to send per check
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

	stats := pinger.Statistics()

	if stats.PacketLoss >= 70.0 {
		result.Message = fmt.Sprintf("FAILED: Not all pings made it back! Received %d out of %d", stats.PacketsRecv, stats.PacketsSent)
		return result
	}

	// TODO: add configuration option to change whether a certain number of pings should make it back, or if a certain percentage of packetloss should be allowed
	// Check for failure of ICMP
	// if stats.PacketsRecv != d.Count {
	// 	result.Message = fmt.Sprintf("FAILED: Not all pings made it back! Received %d out of %d", stats.PacketsRecv, stats.PacketsSent)
	// 	return result
	// }

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
