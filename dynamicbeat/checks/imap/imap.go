package imap

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the imap check
// it implements the "check" interface
type Definition struct {
	ID          string  // a unique identifier for this check
	Name        string  // a human-readable title for the check
	Group       string  // the group this check is part of
	ScoreWeight float64 // the weight that this check has relative to others
	Host        string  // (required) IP or hostname for the imap server
	Username    string  // (required) Username for the imap server
	Password    string  // (required) Password for the user of the imap server
	Encrypted   bool    // (optional, default=false) Whether or not to use TLS (IMAPS)
	Port        string  // (optional, default=143) Port for the imap server
}

// Run a single instance of the check
// We are only supporting the listing of mailboxes as a check currently
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "imap",
	}

	// Create a dialer so we can set timeouts
	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Defining these allow the if/else block below
	var c *client.Client
	var err error

	// Connect to server with TLS or not
	if d.Encrypted {
		c, err = client.DialWithDialerTLS(&dialer, fmt.Sprintf("%s:%s", d.Host, d.Port), &tls.Config{})
	} else {
		c, err = client.DialWithDialer(&dialer, fmt.Sprintf("%s:%s", d.Host, d.Port))
	}
	if err != nil {
		result.Message = fmt.Sprintf("Connecting to server %s failed : %s", d.Host, err)
		out <- result
		return
	}
	defer c.Logout() // This is the same as close() for normal conn objects

	// Set timeout for commands
	c.Timeout = 5 * time.Second

	// Login
	err = c.Login(d.Username, d.Password)
	if err != nil {
		result.Message = fmt.Sprintf("Login with user %s failed : %s", d.Username, err)
		out <- result
		return
	}

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	err = c.List("", "*", mailboxes)
	if err != nil {
		result.Message = fmt.Sprintf("Listing mailboxes failed : %s", err)
		out <- result
		return
	}

	// If we make it here the check passes
	result.Passed = true
	out <- result
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, scoreWeight float64, def []byte) error {

	// Set optional values
	d.Port = "143"

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
	if d.Host == "" {
		missingFields = append(missingFields, "Host")
	}

	if d.Username == "" {
		missingFields = append(missingFields, "Username")
	}

	if d.Password == "" {
		missingFields = append(missingFields, "Password")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "imap",
			Field: missingFields[0],
		}
	}
	return nil
}
