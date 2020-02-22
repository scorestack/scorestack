package winrm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/logp"
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
func (d *Definition) Run(ctx context.Context) schema.CheckResult {

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "winrm",
	}

	// Convert d.Port to int
	port, err := strconv.Atoi(d.Port)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to convert d.Port to int : %s", err)
		return result
	}

	// Another timeout for the bois
	params := *winrm.DefaultParameters
	params.Timeout = "15"

	// Login to winrm and create client
	// endpoint := winrm.NewEndpoint(d.Host, port, d.Encrypted, true, nil, nil, nil, 5*time.Second)
	endpoint := winrm.NewEndpoint(d.Host, port, d.Encrypted, true, nil, nil, nil, 20*time.Second)
	client, err := winrm.NewClientWithParameters(endpoint, d.Username, d.Password, &params)
	if err != nil {
		result.Message = fmt.Sprintf("Login to WinRM host %s failed : %s", d.Host, err)
		return result
	}

	shell, err := client.CreateShell()
	if err != nil {
		result.Message = fmt.Sprintf("Failed to create shell : %s", err)
		return result
	}
	defer func() {
		if closeErr := shell.Close(); closeErr != nil {
			logp.Warn("failed to close winrm connection: %s", closeErr.Error())
		}
	}()

	cmd, err := shell.Execute(d.Cmd)
	if err != nil {
		result.Message = fmt.Sprintf("Executing command %s failed : %s", d.Cmd, err)
		return result
	}

	var test sync.WaitGroup
	copyFunc := func(w io.Writer, r io.Reader) {
		defer test.Done()
		io.Copy(w, r)
	}

	bufOut := new(bytes.Buffer)

	if cmd.Stdout != nil {
		test.Add(1)
		go copyFunc(bufOut, cmd.Stdout)
	} else {
		result.Message = fmt.Sprintf("Failed to get stdout from command %s : %s", d.Cmd, err)
		return result
	}

	cmd.Wait()
	test.Wait()

	// command := winrm.Powershell(d.Cmd)

	// // shell := client.NewShell("ScoreStack-Shell-ID")
	// // defer shell.Close()

	// // cmdOut, err := shell.Execute(command)
	// // defer cmdOut.Close()

	// // if err != nil {
	// // 	result.Message = fmt.Sprintf("Command %s failed : %s", d.Cmd, err)
	// // 	failed <- true
	// // 	return
	// // }

	// // Define these for the command output
	// bufOut := new(bytes.Buffer)
	// bufErr := new(bytes.Buffer)

	// _, err = client.Run(command, bufOut, bufErr)
	// if err != nil {
	// 	result.Message = fmt.Sprintf("Running command %s failed : %s", d.Cmd, err)
	// 	failed <- true
	// 	return
	// }

	// // Check if the command errored
	// if bufErr.String() != "" {
	// 	result.Message = fmt.Sprintf("Executing command %s failed : %s", d.Cmd, bufErr.String())
	// 	failed <- true
	// 	return
	// }

	// // Check if we matching content and the command did not error
	// if !d.MatchContent {
	// 	// If we make it here, no content matching, the check succeeds
	// 	result.Message = fmt.Sprintf("Command %s executed seccessfully: %s", d.Cmd, bufOut.String())
	// 	done <- true
	// 	return
	// }

	// // Keep going if we are matching content
	// // Create regexp
	// regex, err := regexp.Compile(d.ContentRegex)
	// if err != nil {
	// 	result.Message = fmt.Sprintf("Error compiling regex string %s : %s", d.ContentRegex, err)
	// 	failed <- true
	// 	return
	// }

	// // Check if the content matches
	// if !regex.Match(bufOut.Bytes()) {
	// 	result.Message = fmt.Sprintf("Matching content not found")
	// 	failed <- true
	// 	return
	// }

	// If we reach here the check is successful
	result.Passed = true
	return result
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
