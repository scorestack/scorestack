package noop

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of a Noop check.
type Definition struct {
	Config  schema.CheckConfig // generic metadata about the check
	Dynamic string             // (required) contains attributes that can be modified by admins or users
	Static  string             // (required) contains no attributes
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

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(config schema.CheckConfig, def []byte) error {
	// Unpack definition json
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set ID and name attributes
	d.Config = config

	// Make sure required fields are defined
	missingFields := make([]string, 0)
	if d.Dynamic == "" {
		missingFields = append(missingFields, "Dynamic")
	}

	if d.Static == "" {
		missingFields = append(missingFields, "Static")
	}

	// Error on only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.Config.ID,
			Type:  "noop",
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

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
