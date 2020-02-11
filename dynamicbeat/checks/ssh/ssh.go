package ssh

import (
	"encoding/json"
	"sync"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SSH check
// it implements the "check" interface
type Definition struct {
	ID       string // unique identifier for this check
	Name     string // a human-readable title for the check
	IP       string // (required) IP of the host to run the ICMP check against
	Username string // (required) The user to login with over ssh
	Password string // (required) The password for the user that you wish to login with
	Cmd      string // (required) The command to execute once ssh connection established
}

func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {

}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, def []byte) error {

	// Set ID and Name
	d.ID = id
	d.Name = name

	// Unpack JSON definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Check for missing fields
	missingFields := make([]string, 0)
	if d.IP == "" {
		missingFields = append(missingFields, "IP")
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
			Type:  "ssh",
			Field: missingFields[0],
		}
	}
	return nil
}
