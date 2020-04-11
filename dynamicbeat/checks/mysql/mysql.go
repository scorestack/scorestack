package mysql

import (
	"context"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the MySQL check
// it implements the "check" interface
type Definition struct {
	Config       schema.CheckConfig // generic metadata about the check
	Host         string             `optiontype:"required"`                      // IP of Hostname for the MySQL server
	Username     string             `optiontype:"required"`                      // Username for the database
	Password     string             `optiontype:"required"`                      // Password for the user
	Database     string             `optiontype:"required"`                      // Name of the database to access
	Table        string             `optiontype:"required"`                      // Name of the table to access
	Column       string             `optiontype:"required"`                      // Name of the column to access
	ContentRegex string             `optiontype:"optional" optiondefault:".*"`   // Regex to match on
	Port         string             `optiontype:"optional" optiondefault:"3306"` // Port for the server
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
