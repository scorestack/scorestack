package schema

import (
	"context"
	"fmt"
)

// A Check represents the configuration required to verify the operation of a
// single network service.
type Check interface {
	GetConfig() CheckConfig
	SetConfig(config CheckConfig)
	Run(ctx context.Context) CheckResult
	Init(config CheckConfig, def []byte) error
}

// CheckConfig contains the generic metadata about a check that all check types
// must have. It also includes the untemplated check-specific definition JSON
// string and any related attributes. The CheckConfig should be stored in
// structs that implement the Check interface. The Definition and Attribs
// members of the CheckConfig struct will not typically be used by implementers
// of the Check interface, but could still be useful for administrative and
// debugging purposes later on.
type CheckConfig struct {
	ID          string
	Name        string
	Type        string
	Group       string
	ScoreWeight float64
	Definition  []byte
	Attribs     map[string]string
}

// A CheckResult contains information on how a check finished - if it passed or
// failed, and if it failed, why it may have failed. Optionally, a check may
// populate the Details map with additional information about the check's
// execution.
type CheckResult struct {
	Passed  bool
	Message string
	Details map[string]string
}

// A ValidationError represents an issue with a check definition.
type ValidationError struct {
	ID    string // the ID of the check with an invalid definition
	Type  string // the type of the check with an invalid definition
	Field string // the field in the check definition that was invalid
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Error: check (Type: `%s`, ID: `%s`) is missing value for required field `%s`", v.Type, v.ID, v.Field)
}
