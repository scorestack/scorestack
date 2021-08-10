package noop

import (
	"context"
	"strings"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
)

// The Definition configures the behavior of a Noop check.
type Definition struct {
	Config  check.Config // generic metadata about the check
	Dynamic string       `optiontype:"required"` // contains attributes that can be modified by admins or users
	Static  string       `optiontype:"required"` // contains no attributes
}

// Run a single instance of the check.
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialize empty result
	result := check.Result{Timestamp: time.Now(), Metadata: d.Config.Metadata}

	// "Run" the check
	result.Passed = true
	result.Message = strings.Join([]string{d.Dynamic, d.Static}, "; ")
	result.Details = map[string]string{
		"Dynamic": d.Dynamic,
		"Static":  d.Static,
	}

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
