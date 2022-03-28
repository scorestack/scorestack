package ssh

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// The Definition configures the behavior of the SSH check
// it implements the "check" interface
type Definition struct {
	Config       models.CheckConfig // generic metadata about the check
	Host         string             `optiontype:"required"`                    // IP or hostname of the host to run the SSH check against
	Username     string             `optiontype:"required"`                    // The user to login with over ssh
	Password     string             `optiontype:"required"`                    // The password for the user that you wish to login with
	Cmd          string             `optiontype:"required"`                    // The command to execute once ssh connection established
	MatchContent string             `optiontype:"optional"`                    // Whether or not to match content like checking files
	ContentRegex string             `optiontype:"optional" optiondefault:".*"` // Regex to match if reading a file
	Port         string             `optiontype:"optional" optiondefault:"22"` // The port to attempt an ssh connection on
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialize empty result
	result := check.Result{Timestamp: time.Now(), CheckMetadata: d.Config.CheckMetadata}

	// Config SSH client
	// TODO: change timeout to be relative to the parent context's timeout
	config := &ssh.ClientConfig{
		User: d.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(d.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}

	// Create the ssh client
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", d.Host, d.Port), config)
	if err != nil {
		result.Message = fmt.Sprintf("Error creating ssh client: %s", err)
		return result
	}
	defer func() {
		err = client.Close()
		if err != nil {
			zap.S().Warnf("Failed to close SSH connection: %s", err)
		}
	}()

	// Create a session from the connection
	session, err := client.NewSession()
	if err != nil {
		result.Message = fmt.Sprintf("Error creating a ssh session: %s", err)
		return result
	}
	defer func() {
		err = session.Close()
		if err != nil && err.Error() != "EOF" {
			zap.S().Warnf("Failed to close SSH session connection: %s", err)
		}
	}()

	// Run a command
	output, err := session.CombinedOutput(d.Cmd)
	if err != nil {
		result.Message = fmt.Sprintf("Error executing command: %s", err)
		return result
	}

	// Check if we are going to match content
	if matchContent, _ := strconv.ParseBool(d.MatchContent); !matchContent {
		// If we made it here the check passes
		result.Message = fmt.Sprintf("Command %s executed successfully: %s", d.Cmd, output)
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
	if !regex.Match(output) {
		result.Message = "Matching content not found"
		return result
	}

	// If we reach here the check is successful
	result.Passed = true
	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() models.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(c models.CheckConfig) {
	d.Config = c
}
