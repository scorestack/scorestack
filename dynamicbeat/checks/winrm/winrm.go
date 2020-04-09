package winrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/oneNutW0nder/winrm"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the WinRM check
// it implements the "check" interface
type Definition struct {
	Config       schema.CheckConfig // generic metadata about the check
	Host         string             // (required) IP or hostname of the WinRM box
	Username     string             // (required) User to login as
	Password     string             // (required) Password for the user
	Cmd          string             // (required) Command that will be executed
	Encrypted    bool               // (optional, default=true) Use TLS for connection
	MatchContent bool               // (optional, default=false) Turn this on to match content from the output of the cmd
	ContentRegex string             // (optional, default=`.*`) Regexp for matching output of a command
	Port         string             // (optional, default=5986) Port for WinRM
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Convert d.Port to int
	port, err := strconv.Atoi(d.Port)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to convert d.Port to int : %s", err)
		return result
	}

	// Another timeout for the bois
	params := *winrm.DefaultParameters

	// Login to winrm and create client
	// endpoint := winrm.NewEndpoint(d.Host, port, d.Encrypted, true, nil, nil, nil, 5*time.Second)
	endpoint := winrm.NewEndpoint(d.Host, port, d.Encrypted, true, nil, nil, nil, 20*time.Second)
	client, err := winrm.NewClientWithParameters(endpoint, d.Username, d.Password, &params)
	if err != nil {
		result.Message = fmt.Sprintf("Login to WinRM host %s failed : %s", d.Host, err)
		return result
	}

	powershellCmd := winrm.Powershell(d.Cmd)

	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	_, err = client.Run(powershellCmd, bufOut, bufErr)
	if err != nil {
		result.Message = fmt.Sprintf("Executing command %s failed : %s", d.Cmd, err)
		return result
	}

	// Check for an error
	if bufErr.String() != "" {
		result.Message = fmt.Sprintf("Command %s failed : %s", d.Cmd, bufErr.String())
		return result
	}

	// Check if we are going to regex
	if !d.MatchContent {
		// If we make it in here the check passes
		result.Passed = true
		return result
	}

	// Match some content
	regex, err := regexp.Compile(d.ContentRegex)
	if err != nil {
		result.Message = fmt.Sprintf("Error compiling regex string %s : %s", d.ContentRegex, err)
		return result
	}

	// Check if the content matches
	if !regex.Match(bufOut.Bytes()) {
		result.Message = fmt.Sprintf("Matching content not found")
		return result
	}

	// If we reach here the check is successful
	result.Passed = true
	return result
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(config schema.CheckConfig, def []byte) error {

	// Explicitly set defaults
	d.Encrypted = true
	d.ContentRegex = ".*"
	d.Port = "5986"

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

	if d.Cmd == "" {
		missingFields = append(missingFields, "Cmd")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.Config.ID,
			Type:  "winrm",
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

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
