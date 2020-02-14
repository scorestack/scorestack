package winrm

import (
	"encoding/json"
	"sync"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the WinRM check
// it implements the "check" interface
type Definition struct {
	ID           string // a unique identifier for this check
	Name         string // a human-readable title for the check
	Group        string // the group this check is part of
	Host         string // (required) IP or hostname of the WinRM box
	Username     string // (required) User to login as
	Password     string // (required) Password for the user
	Encrypted    bool   // (optional, default=true) Use TLS for connection
	Cmd          string // (required) Command that will be executed
	MatchContent bool   // (optional, default=false) Turn this on to match content from the output of the cmd
	ContentRegex string // (optional, default=`.*`) Regexp for matching output of a command
	Port         string // (optional, default=5869) Port for WinRM
}

// Run a single instance of the check
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {

}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, def []byte) error {

	// Explicitly set defaults
	d.Encrypted = true
	d.ContentRegex = ".*"
	d.Port = "5869"

	// Unpack JSON definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set generic values
	d.ID = id
	d.Name = name
	d.Group = group

	// Check for missing fields
	missingFields := make([]string, 0)
	if d.Host == "" {
		missingFields = append(missingFields, "Host")
	}

	if d.Username == "" {
		missingFields = append(missingFields, "Username")
	}

	if d.Password == "" {
		missingFields = append(missingFields, "Password")
	}

	if d.Cmd == "" {
		missingFields = append(missingFields, "Cmd")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "winrm",
			Field: missingFields[0],
		}
	}
	return nil
}
