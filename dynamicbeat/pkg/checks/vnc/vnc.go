package vnc

import (
	"context"
	"fmt"
	"net"

	vnc "github.com/mitchellh/go-vnc"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/schema"
)

// The Definition configures the behavior of the VNC check
// it implements the "check" interface
type Definition struct {
	Config   schema.CheckConfig // generic metadata about the check
	Host     string             `optiontype:"required"` // The IP or hostname of the vnc server
	Port     string             `optiontype:"required"` // The port for the vnc server
	Password string             `optiontype:"required"` // The password for the vnc server
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Configure the vnc client
	config := vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{Password: d.Password},
		},
	}

	// Make a dialer
	dialer := net.Dialer{}

	// Dial the vnc server
	// conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", d.Host, d.Port), 5*time.Second)
	// TODO: create child context with deadline less than the parent context
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", d.Host, d.Port))
	if err != nil {
		result.Message = fmt.Sprintf("Connection to VNC host %s failed : %s", d.Host, err)
		return result
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			// logp.Warn("Failed to close VNC connection: %s", err)
		}
	}()

	vncClient, err := vnc.Client(conn, &config)
	if err != nil {
		result.Message = fmt.Sprintf("Login to server %s failed : %s", d.Host, err)
		return result
	}
	defer func() {
		err = vncClient.Close()
		if err != nil {
			// logp.Warn("Failed to close VNC connection: %s", err)
		}
	}()

	// If we made it here the check passes
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
