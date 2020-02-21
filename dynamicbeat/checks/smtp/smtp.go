package smtp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/smtp"
	"sync"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SMTP check
// it implements the "check" interface
type Definition struct {
	ID          string  // a unique identifier for this check
	Name        string  // a human-readable title for the check
	Group       string  // the group this check is part of
	ScoreWeight float64 // the weight that this check has relative to others
	Host        string  // (required) IP or hostname of the smtp server
	Username    string  // (required) Username for the smtp server
	Password    string  // (required) Password for the smtp server
	Sender      string  // (required) Who is sending the email
	Reciever    string  // (required) Who is receiving the email
	Body        string  // (optional, default="Hello from ScoreStack") Body of the email
	Encrypted   bool    // (optional, default=false) Whether or not to use TLS
	Port        string  // (optional, default="25") Port of the smtp server
}

// **************************************************
type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

// **************************************************

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
		CheckType:   "smtp",
	}

	// Make channels for ocmpleting the check or not
	done := make(chan bool)
	failed := make(chan bool)

	go func() {
		// Create a dialer
		dialer := net.Dialer{
			Timeout: 5 * time.Second,
		}

		// ***********************************************
		// Set up custom auth for bypassing net/smtp protections
		auth := unencryptedAuth{smtp.PlainAuth("", d.Username, d.Password, d.Host)}
		// ***********************************************

		// The good way to do auth
		// auth := smtp.PlainAuth("", d.Username, d.Password, d.Host)
		// Create TLS config
		tlsConfig := tls.Config{
			InsecureSkipVerify: true,
		}

		// Declare these for the below if block
		var conn net.Conn
		var err error

		if d.Encrypted {
			conn, err = tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:%s", d.Host, d.Port), &tlsConfig)
		} else {
			conn, err = dialer.Dial("tcp", fmt.Sprintf("%s:%s", d.Host, d.Port))
		}
		if err != nil {
			result.Message = fmt.Sprintf("Connecting to server %s failed : %s", d.Host, err)
			failed <- true
			return
		}
		defer conn.Close()

		// Create smtp client
		c, err := smtp.NewClient(conn, d.Host)
		if err != nil {
			result.Message = fmt.Sprintf("Created smtp client to host %s failed : %s", d.Host, err)
			failed <- true
			return
		}
		defer c.Quit()

		// Login
		err = c.Auth(auth)
		if err != nil {
			result.Message = fmt.Sprintf("Login to %s failed : %s", d.Host, err)
			failed <- true
			return
		}

		// Set the sender
		err = c.Mail(d.Sender)
		if err != nil {
			result.Message = fmt.Sprintf("Setting sender %s failed : %s", d.Sender, err)
			failed <- true
			return
		}

		// Set the reciver
		err = c.Rcpt(d.Reciever)
		if err != nil {
			result.Message = fmt.Sprintf("Setting reciever %s failed : %s", d.Reciever, err)
			failed <- true
			return
		}

		// Send the email body.
		wc, err := c.Data()
		if err != nil {
			result.Message = fmt.Sprintf("Creating writer failed : %s", err)
			failed <- true
			return
		}
		defer wc.Close()

		// Write the body
		_, err = fmt.Fprintf(wc, d.Body)
		if err != nil {
			result.Message = fmt.Sprintf("Writing mail body failed : %s", err)
			failed <- true
			return
		}

		// If we make it here the check succeeds
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

	// Explicitly set defaults
	d.Body = "Hello from ScoreStack"
	d.Port = "25"

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

	if d.Sender == "" {
		missingFields = append(missingFields, "Sender")
	}

	if d.Reciever == "" {
		missingFields = append(missingFields, "Reciever")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "smtp",
			Field: missingFields[0],
		}
	}
	return nil
}
