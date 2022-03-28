package noop

import (
	"context"
	"strings"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/models"
)

// The Definition configures the behavior of a Noop check.
type Definition struct {
	Config  models.CheckConfig // generic metadata about the check
	Dynamic string             `optiontype:"required"` // contains attributes that can be modified by admins or users
	Static  string             `optiontype:"required"` // contains no attributes
}

// Run a single instance of the check.
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialize empty result
	result := check.Result{Timestamp: time.Now(), CheckMetadata: d.Config.CheckMetadata}

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
func (d *Definition) GetConfig() models.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(c models.CheckConfig) {
	d.Config = c
}
