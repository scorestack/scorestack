package ldap

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
	"gopkg.in/ldap.v2"
)

// The Definition configures the behavior of the LDAP check
// it implements the "check" interface
type Definition struct {
	ID          string  // a unique identifier for this check
	Name        string  // a human-readable title for the check
	Group       string  // the group this check is part of
	ScoreWeight float64 // the weight that this check has relative to others
	User        string  // (required) The user written in DN syntax
	Password    string  // (required) the password for the user
	Fqdn        string  // (required) The Fqdn of the ldap server
	Ldaps       bool    // (optional, default=false) Whether or not to use LDAP+TLS
	Port        string  // (optional, default=389) Port for ldap
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context, wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "ldap",
	}

	// Make channels for completing the check or not
	done := make(chan bool)
	failed := make(chan bool)

	go func() {
		// Set timeout
		ldap.DefaultTimeout = 5 * time.Second

		// Normal, default ldap check
		lconn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", d.Fqdn, d.Port))
		if err != nil {
			result.Message = fmt.Sprintf("Could not dial server %s : %s", d.Fqdn, err)
			failed <- true
			return
		}
		defer lconn.Close()

		// Set message timeout
		lconn.SetTimeout(5 * time.Second)

		// Add TLS if needed
		if d.Ldaps {
			err = lconn.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				result.Message = fmt.Sprintf("TLS session creation failed : %s", err)
				failed <- true
				return
			}
		}

		// Attempt to login
		err = lconn.Bind(d.User, d.Password)
		if err != nil {
			result.Message = fmt.Sprintf("Failed to login with user %s : %s", d.User, err)
			failed <- true
			return
		}

		// If we reached here the check passes
		done <- true
	}()

	// Watch channels and context for timeout
	for {
		select {
		case <-done:
			result.Passed = true
			out <- result
			return
		case <-failed:
			out <- result
			return
		case <-ctx.Done():
			result.Message = fmt.Sprintf("Timeout via context : %s", ctx.Err())
			out <- result
			return
		}
	}
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, scoreWeight float64, def []byte) error {

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
	d.ScoreWeight = scoreWeight

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
