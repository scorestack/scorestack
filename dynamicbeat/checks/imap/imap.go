package imap

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the imap check
// it implements the "check" interface
type Definition struct {
	Config    schema.CheckConfig // generic metadata about the check
	Host      string             `optiontype:"required"`                     // IP or hostname for the imap server
	Username  string             `optiontype:"required"`                     // Username for the imap server
	Password  string             `optiontype:"required"`                     // Password for the user of the imap server
	Encrypted bool               `optiontype:"optional"`                     // Whether or not to use TLS (IMAPS)
	Port      string             `optiontype:"optional" optiondefault:"143"` // Port for the imap server
}

// Run a single instance of the check
// We are only supporting the listing of mailboxes as a check currently
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Create a dialer so we can set timeouts
	// TODO: change this to be relative to the parent context's timeout
	dialer := net.Dialer{
		Timeout: 20 * time.Second,
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
		return result
	}
	defer func() {
		if closeErr := c.Logout(); closeErr != nil {
			// logp.Warn("failed to close imap connection: %s", closeErr.Error())
		}
	}()

	// Set timeout for commands
	c.Timeout = 5 * time.Second

	// Login
	err = c.Login(d.Username, d.Password)
	if err != nil {
		result.Message = fmt.Sprintf("Login with user %s failed : %s", d.Username, err)
		return result
	}

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	err = c.List("", "*", mailboxes)
	if err != nil {
		result.Message = fmt.Sprintf("Listing mailboxes failed : %s", err)
		return result
	}

	// If we make it here the check passes
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
