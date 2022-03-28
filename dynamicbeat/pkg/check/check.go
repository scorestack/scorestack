package check

import (
	"context"
	"fmt"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/models"
)

type Check interface {
	GetConfig() models.CheckConfig
	SetConfig(c models.CheckConfig)
	Run(ctx context.Context) Result
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
