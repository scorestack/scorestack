package ftp

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

// The Definition configures the behavior of the FTP check
// it implements the "check" interface
type Definition struct {
	Config           check.Config // generic metadata about the check
	Host             string       `optiontype:"required"`                    // IP or hostname of the host to run the FTP check against
	Username         string       `optiontype:"required"`                    // The user to login with over FTP
	Password         string       `optiontype:"required"`                    // The password for the user that you wish to login with
	File             string       `optiontype:"required"`                    // The path to the file to access during the FTP check
	ContentRegex     string       `optiontype:"optional" optiondefault:".*"` // Regex to match if reading a file
	HashContentMatch string       `optiontype:"optional"`                    // Whether or not to match a hash of the file contents
	Hash             string       `optiontype:"optional"`                    // The hash digest from sha3-256 to compare the hashed file contents to
	Port             string       `optiontype:"optional" optiondefault:"21"` // The port to attempt an ftp connection on
	Simple           string       `optiontype:"optional"`                    // Very simple FTP check for older servers
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialize empty result
	result := check.Result{Timestamp: time.Now(), Metadata: d.Config.Metadata}

	// Connect to the ftp server
	// TODO: create child context with deadline less than the parent context
	conn, err := ftp.Dial(fmt.Sprintf("%s:%s", d.Host, d.Port), ftp.DialWithContext(ctx))
	if err != nil {
		result.Message = fmt.Sprintf("Connection to %s on port %s failed : %s", d.Host, d.Port, err)
		return result
	}
	defer func() {
		err := conn.Quit()
		if err != nil {
			zap.S().Warnf("Failed to close FTP connection: %s", err)
		}
	}()

	// Login
	err = conn.Login(d.Username, d.Password)
	if err != nil {
		result.Message = fmt.Sprintf("Login attempt with user %s failed : %s", d.Username, err)
		return result
	}

	// ***********************************************
	if simple, _ := strconv.ParseBool(d.Simple); simple {
		// Do a simple FTP check for servers that don't support a lot of FTP commands
		err = conn.ChangeDir(d.File)
		if err != nil {
			result.Message = fmt.Sprintf("Changing to directory %s failed : %s", d.File, err)
			return result
		}

		_, err := conn.CurrentDir()
		// entries, err := conn.List("/")
		if err != nil {
			result.Message = fmt.Sprintf("Getting current directory %s failed : %s", d.File, err)
			return result
		}

		// If we reached here, changed dir success, check passed
		result.Passed = true
		return result
	}
	// **************

	// Retrieve file contents
	resp, err := conn.Retr(d.File)
	if err != nil {
		result.Message = fmt.Sprintf("Could not retrieve file %s : %s", d.File, err)
		return result
	}
	defer resp.Close()

	content, err := io.ReadAll(resp)
	if err != nil {
		result.Message = fmt.Sprintf("Could not read file %s contents : %s", d.File, err)
		return result
	}

	// Check if we are doing hash matching, non default
	if matchHash, _ := strconv.ParseBool(d.HashContentMatch); matchHash {
		// Get the file hash
		digest := sha3.Sum256(content)

		// Check if the digest of the file matches the defined hash
		if digestString := hex.EncodeToString(digest[:]); digestString != d.Hash {
			result.Message = "Incorrect hash"
			return result
		}

		// If we make it here the check was successful for matching hashes
		result.Passed = true
		return result
	}

	// Default, regex content matching
	regex, err := regexp.Compile(d.ContentRegex)
	if err != nil {
		result.Message = fmt.Sprintf("Error compiling regex string %s : %s", d.ContentRegex, err)
		return result
	}

	// Check if content matches regex
	if !regex.Match(content) {
		result.Message = "Matching content not found"
		return result
	}

	// If we reach here the check is successful
	result.Passed = true
	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with.
func (d *Definition) GetConfig() check.Config {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(c check.Config) {
	d.Config = c
}
