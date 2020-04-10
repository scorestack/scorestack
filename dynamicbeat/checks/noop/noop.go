package noop

import (
	"context"
	"strings"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of a Noop check.
type Definition struct {
	Config  schema.CheckConfig // generic metadata about the check
	Dynamic string             `optiontype:"required"` // contains attributes that can be modified by admins or users
	Static  string             `optiontype:"required"` // contains no attributes
}

// Run a single instance of the check.
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

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
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
