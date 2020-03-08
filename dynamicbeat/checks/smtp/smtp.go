package smtp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SMTP check
// it implements the "check" interface
type Definition struct {
	Config    schema.CheckConfig // generic metadata about the check
	Host      string             // (required) IP or hostname of the smtp server
	Username  string             // (required) Username for the smtp server
	Password  string             // (required) Password for the smtp server
	Sender    string             // (required) Who is sending the email
	Reciever  string             // (required) Who is receiving the email
	Body      string             // (optional, default="Hello from ScoreStack") Body of the email
	Encrypted bool               // (optional, default=false) Whether or not to use TLS
	Port      string             // (optional, default="25") Port of the smtp server
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
func (d *Definition) Run(ctx context.Context, result schema.CheckResult) schema.CheckResult {

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.Config.ID,
		Name:        d.Config.Name,
		Group:       d.Config.Group,
		ScoreWeight: d.Config.ScoreWeight,
		CheckType:   "smtp",
	}

	// Create a dialer
	dialer := net.Dialer{
		Timeout: 20 * time.Second,
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

	// KILL THIS CHECK!
	done := make(chan bool)

	go func() {
		// Declare these for the below if block
		var conn net.Conn
		var err error

		if d.Encrypted {
			conn, err = tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:%s", d.Host, d.Port), &tlsConfig)
		} else {
			conn, err = dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", d.Host, d.Port))
		}
		if err != nil {
			result.Message = fmt.Sprintf("Connecting to server %s failed : %s", d.Host, err)
			done <- true
			return
		}
		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				// logp.Warn("failed to close smtp connection: %s", closeErr.Error())
			}
		}()

		// Create smtp client
		c, err := smtp.NewClient(conn, d.Host)
		if err != nil {
			result.Message = fmt.Sprintf("Created smtp client to host %s failed : %s", d.Host, err)
			done <- true
			return
		}
		defer func() {
			if closeErr := c.Quit(); closeErr != nil {
				// logp.Warn("failed to close smtp client connection: %s", closeErr.Error())
			}
		}()

		// Login
		err = c.Auth(auth)
		if err != nil {
			result.Message = fmt.Sprintf("Login to %s failed : %s", d.Host, err)
			done <- true
			return
		}

		// Set the sender
		err = c.Mail(d.Sender)
		if err != nil {
			result.Message = fmt.Sprintf("Setting sender %s failed : %s", d.Sender, err)
			done <- true
			return
		}

		// Set the reciver
		err = c.Rcpt(d.Reciever)
		if err != nil {
			result.Message = fmt.Sprintf("Setting reciever %s failed : %s", d.Reciever, err)
			done <- true
			return
		}

		// Send the email body.
		wc, err := c.Data()
		if err != nil {
			result.Message = fmt.Sprintf("Creating writer failed : %s", err)
			done <- true
			return
		}
		defer wc.Close()

		// Write the body
		_, err = fmt.Fprintf(wc, d.Body)
		if err != nil {
			result.Message = fmt.Sprintf("Writing mail body failed : %s", err)
			done <- true
			return
		}

		result.Passed = true
	}()

	for {
		select {
		case <-ctx.Done():
			result.Message = fmt.Sprintf("Timeout limit reached: %s", ctx.Err())
			return result
		case <-done:
			return result
		}
	}
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(config schema.CheckConfig, def []byte) error {

	// Explicitly set defaults
	d.Body = "Hello from ScoreStack"
	d.Port = "25"

	// Unpack JSON definition
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set generic values
	d.Config = config

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
			ID:    d.Config.ID,
			Type:  "smtp",
			Field: missingFields[0],
		}
	}
	return nil
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}
