package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/scorestack/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SMTP check
// it implements the "check" interface
type Definition struct {
	Config    schema.CheckConfig // generic metadata about the check
	Host      string             `optiontype:"required"`                                       // IP or hostname of the smtp server
	Username  string             `optiontype:"required"`                                       // Username for the smtp server
	Password  string             `optiontype:"required"`                                       // Password for the smtp server
	Sender    string             `optiontype:"required"`                                       // Who is sending the email
	Reciever  string             `optiontype:"required"`                                       // Who is receiving the email
	Body      string             `optiontype:"optional" optiondefault:"Hello from Scorestack"` // Body of the email
	Encrypted string             `optiontype:"optional"`                                       // Whether or not to use TLS
	Port      string             `optiontype:"optional" optiondefault:"25"`                    // Port of the smtp server
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
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Create a dialer
	// TODO: change this to be relative to the parent context's timeout
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

	// Declare these for the below if block
	var conn net.Conn
	var err error

	if encrypted, _ := strconv.ParseBool(d.Encrypted); encrypted {
		conn, err = tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:%s", d.Host, d.Port), &tlsConfig)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", d.Host, d.Port))
	}
	if err != nil {
		result.Message = fmt.Sprintf("Connecting to server %s failed : %s", d.Host, err)
		return result
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logp.Warn("Failed to close SMTP connection: %s", err)
		}
	}()

	// Create smtp client
	c, err := smtp.NewClient(conn, d.Host)
	if err != nil {
		result.Message = fmt.Sprintf("Created smtp client to host %s failed : %s", d.Host, err)
		return result
	}
	defer func() {
		err = c.Quit()
		if err != nil {
			logp.Warn("Failed to close SMTP client connection: %s", err)
		}
	}()

	// Login
	err = c.Auth(auth)
	if err != nil {
		result.Message = fmt.Sprintf("Login to %s failed : %s", d.Host, err)
		return result
	}

	// Set the sender
	err = c.Mail(d.Sender)
	if err != nil {
		result.Message = fmt.Sprintf("Setting sender %s failed : %s", d.Sender, err)
		return result
	}

	// Set the reciver
	err = c.Rcpt(d.Reciever)
	if err != nil {
		result.Message = fmt.Sprintf("Setting reciever %s failed : %s", d.Reciever, err)
		return result
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		result.Message = fmt.Sprintf("Creating writer failed : %s", err)
		return result
	}
	defer wc.Close()

	// Write the body
	_, err = fmt.Fprintf(wc, d.Body)
	if err != nil {
		result.Message = fmt.Sprintf("Writing mail body failed : %s", err)
		return result
	}

	result.Passed = true
	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
