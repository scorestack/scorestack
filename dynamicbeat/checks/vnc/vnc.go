package vnc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/mitchellh/go-vnc"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the VNC check
// it implements the "check" interface
type Definition struct {
	Config   schema.CheckConfig // generic metadata about the check
	Host     string             // (required) The IP or hostname of the vnc server
	Port     string             // (required) The port for the vnc server
	Password string             // (required) The password for the vnc server
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context, result schema.CheckResult) schema.CheckResult {

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.Config.ID,
		Name:        d.Config.Name,
		Group:       d.Config.Group,
		ScoreWeight: d.Config.ScoreWeight,
		CheckType:   "vnc",
	}

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
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", d.Host, d.Port))
	if err != nil {
		result.Message = fmt.Sprintf("Connection to VNC host %s failed : %s", d.Host, err)
		return result
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			// logp.Warn("failed to close vnc connection: %s", closeErr.Error())
		}
	}()

	vncClient, err := vnc.Client(conn, &config)
	if err != nil {
		result.Message = fmt.Sprintf("Login to server %s failed : %s", d.Host, err)
		return result
	}
	defer func() {
		if closeErr := vncClient.Close(); closeErr != nil {
			// logp.Warn("failed to close vnc connection: %s", closeErr.Error())
		}
	}()

	// If we made it here the check passes
	result.Passed = true
	return result

}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(config schema.CheckConfig, def []byte) error {

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

	if d.Port == "" {
		missingFields = append(missingFields, "Port")
	}

	if d.Password == "" {
		missingFields = append(missingFields, "Password")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.Config.ID,
			Type:  "vnc",
			Field: missingFields[0],
		}
	}
	return nil
}
