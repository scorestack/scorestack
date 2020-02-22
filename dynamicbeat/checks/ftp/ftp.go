package ftp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
	"golang.org/x/crypto/sha3"
)

// The Definition configures the behavior of the FTP check
// it implements the "check" interface
type Definition struct {
	ID                string  // a unique identifier for this check
	Name              string  // a human-readable title for the check
	Group             string  // the group this check is part of
	ScoreWeight       float64 // the weight that this check has relative to others
	Host              string  // (required) IP or hostname of the host to run the FTP check against
	Username          string  // (required) The user to login with over FTP
	Password          string  // (required) The password for the user that you wish to login with
	File              string  // (required) The path to the file to access during the FTP check
	RegexContentMatch bool    // (optional, default=true) Whether or not to match file content with regex
	ContentRegex      string  // (optional, default=`.*`) Regex to match if reading a file
	HashContentMatch  bool    // (optional, default=false) Whether or not to match a hash of the file contents
	Hash              string  // (optional, default="") The hash digest from sha3-256 to compare the hashed file contents to
	Port              string  // (optional, default=21) The port to attempt an ftp connection on
	Fucked            bool    // (optional, default=false) Custom case for Cerealkiller ISTS2020
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context, wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	// Setup result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "ftp",
	}

	// Create channels for passing or failing
	done := make(chan bool)
	failed := make(chan bool)

	go func() {
		// Connect to the ftp server
		// conn, err := ftp.Dial(fmt.Sprintf("%s:%s", d.Host, d.Port), ftp.DialWithTimeout(5*time.Second))
		conn, err := ftp.Dial(fmt.Sprintf("%s:%s", d.Host, d.Port))
		if err != nil {
			result.Message = fmt.Sprintf("Connection to %s on port %s failed : %s", d.Host, d.Port, err)
			failed <- true
			return
		}
		defer conn.Quit()

		// Login
		err = conn.Login(d.Username, d.Password)
		if err != nil {
			result.Message = fmt.Sprintf("Login attempt with user %s failed : %s", d.Username, err)
			failed <- true
			return
		}

		// ***********************************************
		if d.Fucked {
			// Do check for cerealkiller
			err = conn.ChangeDir(d.File)
			if err != nil {
				result.Message = fmt.Sprintf("Changing to directory %s failed : %s", d.File, err)
				failed <- true
				return
			}

			_, err := conn.CurrentDir()
			// entries, err := conn.List("/")
			if err != nil {
				result.Message = fmt.Sprintf("Getting current directory %s failed : %s", d.File, err)
				failed <- true
				return
			}

			// If we reached here, changed dir success, check passed
			done <- true
			return
		}
		// **************

		// Retrieve file contents
		resp, err := conn.Retr(d.File)
		if err != nil {
			result.Message = fmt.Sprintf("Could not retrieve file %s : %s", d.File, err)
			failed <- true
			return
		}
		defer resp.Close()

		content, err := ioutil.ReadAll(resp)
		if err != nil {
			result.Message = fmt.Sprintf("Could not read file %s contents : %s", d.File, err)
			failed <- true
			return
		}

		// Check if we are doing hash matching, non default
		if d.HashContentMatch {
			// Get the file hash
			digest := sha3.Sum256(content)

			// Check if the digest of the file matches the defined hash
			if digestString := hex.EncodeToString(digest[:]); digestString != d.Hash {
				result.Message = fmt.Sprintf("Incorrect hash")
				failed <- true
				return
			}

			// If we make it here the check was successful for matching hashes
			done <- true
			return
		}

		// Default, regex content matching
		regex, err := regexp.Compile(d.ContentRegex)
		if err != nil {
			result.Message = fmt.Sprintf("Error compiling regex string %s : %s", d.ContentRegex, err)
			failed <- true
			return
		}

		// Check if content matches regex
		if !regex.Match(content) {
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

	// Explicitly set default, optional values
	d.RegexContentMatch = true
	d.ContentRegex = ".*"
	d.Port = "21"

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

	if d.File == "" {
		missingFields = append(missingFields, "File")
	}

	// Error only the first missing field, if there are any
	if len(missingFields) > 0 {
		return schema.ValidationError{
			ID:    d.ID,
			Type:  "ftp",
			Field: missingFields[0],
		}
	}
	return nil
}
