package schema

import (
	"fmt"
	"sync"
	"time"
)

// A Check represents the configuration required to verify the operation of a
// single network service.
type Check interface {
	Run(wg *sync.WaitGroup, out chan<- CheckResult)
	Init(id string, name string, def []byte) error
}

// A CheckDef is an untemplated representation of a check. In this format, the
// definition is represented as a JSON string.
type CheckDef struct {
	ID         string
	Name       string
	Type       string
	Definition []byte
	Attribs    map[string]string
}

// CheckResult : Information about the results of executing a check.
type CheckResult struct {
	Timestamp time.Time
	ID        string
	Name      string
	CheckType string
	Passed    bool
	Message   string
	Details   map[string]string
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
