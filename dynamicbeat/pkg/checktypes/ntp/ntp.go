package ntp

import (
	"context"
	"fmt"
	"time"

	"github.com/beevik/ntp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
)

// The Definition configures the behavior of the NTP check
// it implements the "check" interface
type Definition struct {
	Config check.Config // generic metadata about the check
	Fqdn   string       `optiontype:"required"` // The FQDN of the host you are looking up
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialize empty result
	result := check.Result{Timestamp: time.Now(), Metadata: d.Config.Metadata}

	_, err := ntp.Query(d.Fqdn)
	if err != nil {
		result.Message = fmt.Sprintf("Problem sending query to %s : %s", d.Fqdn, err)
		return result
	} else {
		result.Passed = true
		return result
	}
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
