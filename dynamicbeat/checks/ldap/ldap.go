package ldap

import (
	"encoding/json"
	"sync"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SSH check
// it implements the "check" interface
type Definition struct {
	ID       string // a unique identifier for this check
	Name     string // a human-readable title for the check
	Group    string // the group this check is part of
	User     string // (required) The user written in DN syntax
	Password string // (required) the password for the user
	Fqdn     string // (required) The Fqdn of the ldap server
	Ldaps    bool   // (optional, default=false) Whether or not to use LDAP+TLS
	Port     string // (optional, default=389) Port for ldap
}

// Run a single instance of the check
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {

}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, def []byte) error {

	// Explicitly set default values
	d.Port = "389"

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
	if d.User == "" {
		missingFields = append(missingFields, "User")
	}

	if d.Password == "" {
		missingFields = append(missingFields, "Password")
	}

	if d.Fqdn == "" {
		missingFields = append(missingFields, "Fqdn")
	}

	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "ldap",
			Field: missingFields[0],
		}
	}
	return nil
}
