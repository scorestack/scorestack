package ldap

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
	"gopkg.in/ldap.v2"
)

// The Definition configures the behavior of the LDAP check
// it implements the "check" interface
type Definition struct {
	Config   schema.CheckConfig // generic metadata about the check
	User     string             `optiontype:"required"`                     // The user written in user@domain syntax
	Password string             `optiontype:"required"`                     // the password for the user
	Fqdn     string             `optiontype:"required"`                     // The Fqdn of the ldap server
	Ldaps    bool               `optiontype:"optional"`                     // Whether or not to use LDAP+TLS
	Port     string             `optiontype:"optional" optiondefault:"389"` // Port for ldap
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Set timeout
	// TODO: change this to be relative to the parent context's timeout
	ldap.DefaultTimeout = 20 * time.Second

	// Normal, default ldap check
	lconn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", d.Fqdn, d.Port))
	if err != nil {
		result.Message = fmt.Sprintf("Could not dial server %s : %s", d.Fqdn, err)
		return result
	}
	defer lconn.Close()

	// Set message timeout
	lconn.SetTimeout(5 * time.Second)

	// Add TLS if needed
	if d.Ldaps {
		err = lconn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			result.Message = fmt.Sprintf("TLS session creation failed : %s", err)
			return result
		}
	}

	// Attempt to login
	err = lconn.Bind(d.User, d.Password)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to login with user %s : %s", d.User, err)
		return result
	}

	// If we reached here the check passes
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
