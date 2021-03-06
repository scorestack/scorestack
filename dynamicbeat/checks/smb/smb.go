package smb

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"regexp"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/hirochachacha/go-smb2"
	"github.com/scorestack/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the SMB check
// it implements the "check" interface
type Definition struct {
	Config       schema.CheckConfig // generic metadata about the check
	Host         string             `optiontype:"required"`                     // IP or hostname for SMB server
	Username     string             `optiontype:"required"`                     // Username for SMB share
	Password     string             `optiontype:"required"`                     // Password for SMB user
	Share        string             `optiontype:"required"`                     // Name of the share
	Domain       string             `optiontype:"required"`                     // The domain found in front of a login (SMB\Administrator : SMB would be the domain)
	File         string             `optiontype:"required"`                     // The file in the SMB share
	ContentRegex string             `optiontype:"optional" optiondefault:".*"`  // Regex to match on
	Port         string             `optiontype:"optional" optiondefault:"445"` // Port of the server
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Dial SMB server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", d.Host, d.Port))
	if err != nil {
		result.Message = fmt.Sprintf("Error with initial dial : %s", err)
		return result
	}
	defer conn.Close()

	// Configure SMB dialer
	smbConn := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     d.Username,
			Password: d.Password,
			Domain:   d.Domain,
		},
	}

	// Dial SMB server for SMB connection
	c, err := smbConn.DialContext(ctx, conn)
	if err != nil {
		result.Message = fmt.Sprintf("Error connecting to smb server : %s", err)
		return result
	}
	defer func() {
		err := c.Logoff()
		if err != nil {
			logp.Info("Error logging off from SMB server: %s", err)
		}
	}()

	// Mount the SMB share
	fs, err := c.Mount(fmt.Sprintf(`\\%s\%s`, d.Host, d.Share))
	if err != nil {
		result.Message = fmt.Sprintf("Error mounting share : %s", err)
		return result
	}
	defer func() {
		err := fs.Umount()
		if err != nil {
			logp.Warn("Error unmounting remote file system: %s", err)
		}
	}()

	// Open the file for reading
	f, err := fs.Open(d.File)
	if err != nil {
		result.Message = fmt.Sprintf("Error opening file : %s", err)
		return result
	}
	defer f.Close()

	// Ensure we are reading from the beginning of the file
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		result.Message = fmt.Sprintf("Error seeking to beginning of file : %s", err)
		return result
	}

	// Read from the file
	content, err := ioutil.ReadAll(f)
	if err != nil {
		result.Message = fmt.Sprintf("Error reading the file contents : %s", err)
		return result
	}

	// Compile regex
	regex, err := regexp.Compile(d.ContentRegex)
	if err != nil {
		result.Message = fmt.Sprintf("Error compiling regex string %s : %s", d.ContentRegex, err)
		return result
	}

	// Check if content matches regex
	if !regex.Match(content) {
		result.Message = fmt.Sprintf("Matching content not found")
		return result
	}

	// If we reach here the check is successful
	result.Passed = true
	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
