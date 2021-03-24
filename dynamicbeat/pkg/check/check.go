package check

import (
	"context"
	"fmt"
	"time"
)

type Check interface {
	GetConfig() Config
	SetConfig(c Config)
	Run(ctx context.Context) Result
}

type Metadata struct {
	ID          string
	Name        string
	Type        string
	Group       string
	ScoreWeight float64
}

type Config struct {
	Metadata
	Definition []byte
	Attribs    map[string]string
}

type Result struct {
	Timestamp time.Time
	Metadata
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
