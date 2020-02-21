package winrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/masterzen/winrm"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the WinRM check
// it implements the "check" interface
type Definition struct {
	ID           string  // a unique identifier for this check
	Name         string  // a human-readable title for the check
	Group        string  // the group this check is part of
	ScoreWeight  float64 // the weight that this check has relative to others
	Host         string  // (required) IP or hostname of the WinRM box
	Username     string  // (required) User to login as
	Password     string  // (required) Password for the user
	Cmd          string  // (required) Command that will be executed
	Encrypted    bool    // (optional, default=true) Use TLS for connection
	MatchContent bool    // (optional, default=false) Turn this on to match content from the output of the cmd
	ContentRegex string  // (optional, default=`.*`) Regexp for matching output of a command
	Port         string  // (optional, default=5986) Port for WinRM
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
		CheckType:   "winrm",
	}

	// make channels for completing the check or not
	done := make(chan bool)
	failed := make(chan bool)

	go func() {
		// Convert d.Port to int
		port, err := strconv.Atoi(d.Port)
		if err != nil {
			result.Message = fmt.Sprintf("Failed to convert d.Port to int : %s", err)
			failed <- true
			return
		}

		// Another timeout for the bois
		params := winrm.Parameters{
			Timeout: "5",
		}

		// Login to winrm and create client
		endpoint := winrm.NewEndpoint(d.Host, port, d.Encrypted, true, nil, nil, nil, 5*time.Second)
		// client, err := winrm.NewClient(endpoint, d.Username, d.Password)
		client, err := winrm.NewClientWithParameters(endpoint, d.Username, d.Password, &params)
		if err != nil {
			result.Message = fmt.Sprintf("Login to WinRM host %s failed : %s", d.Host, err)
			failed <- true
			return
		}
		client.Timeout = "5"

		// Define these for the command output
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)

		// Execute a command
		_, err = client.Run("netstat", bufOut, bufErr)
		if err != nil {
			result.Message = fmt.Sprintf("Running command %s failed : %s", d.Cmd, err)
			failed <- true
			return
		}

		// Check if the command errored
		if bufErr.String() != "" {
			result.Message = fmt.Sprintf("Executing command %s failed : %s", d.Cmd, bufErr.String())
			failed <- true
			return
		}

		// Check if we matching content and the command did not error
		if !d.MatchContent {
			// If we make it here, no content matching, the check succeeds
			result.Message = fmt.Sprintf("Command %s executed seccessfully: %s", d.Cmd, bufOut.String())
			done <- true
			return
		}

		// Keep going if we are matching content
		// Create regexp
		regex, err := regexp.Compile(d.ContentRegex)
		if err != nil {
			result.Message = fmt.Sprintf("Error compiling regex string %s : %s", d.ContentRegex, err)
			failed <- true
			return
		}

		// Check if the content matches
		if !regex.Match(bufOut.Bytes()) {
			result.Message = fmt.Sprintf("Matching content not found")
			failed <- true
			return
		}

		// If we reach here the check is successful
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
	d.Encrypted = true
	d.ContentRegex = ".*"
	d.Port = "5986"

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

	if d.Cmd == "" {
		missingFields = append(missingFields, "Cmd")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "winrm",
			Field: missingFields[0],
		}
	}
	return nil
}
