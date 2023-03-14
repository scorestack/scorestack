package git

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
)

// The Definition configures the behavior of the Git check and implements the "check" interface.
type Definition struct {
	Config          check.Config // Generic metadata about the check
	Host            string       `optiontype:"required"`                    // IIP or FQDN the remote repository is located
	Repository      string       `optiontype:"required"`                    // The path to the remote repository
	Branch          string       `optiontype:"required"`                    // The branch to clone from the repository
	Port            int          `optiontype:"optional"`                    // The port to connect to for cloning the repository
	HTTPS           bool         `optiontype:"optional"`                    // Whether to use HTTP or HTTPS
	HttpsValidate   bool         `optiontype:"optional"`                    // Whether HTTPS certificates should be validated
	Username        string       `optiontype:"optional"`                    // Username to use for private repositories
	Password        string       `optiontype:"optional"`                    // Password for the user
	ContentMatch    bool         `optiontype:"optional"`                    // Whether to check the contents of a file
	ContentFile     string       `optiontype:"optional"`                    // The path of the file to check the contents of
	ContentRegex    string       `optiontype:"optional" optiondefault:".*"` // The regex to match against the checked file
	CommitHashMatch bool         `optiontype:"optional"`                    // Whether or not to match the hash of the latest commit
	CommitHash      string       `optiontype:"optional"`                    // The hash to check against the latest commit
}

// Run a single instance of the check.
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialilze empty result
	result := check.Result{Timestamp: time.Now(), Metadata: d.Config.Metadata}

	// Fail the check if an invalid method is attempted
	var protocol string
	switch d.HTTPS {
	case true:
		protocol = "https"
	case false:
		protocol = "http"
	}

	// Craft the repository url for git.Clone
	// Using the implementation provided by go-git as it's probably better than anything I can come up with
	repoUrl := transport.Endpoint{
		Protocol: protocol,
		Host:     d.Host,
		Port:     d.Port,
		User:     d.Username,
		Password: d.Password,
		Path:     d.Repository,
	}

	// Intialize memory storage for cloning
	store := memory.NewStorage()
	tree := memfs.New()

	// Clone the Git repository into memory
	// We only clone the latest commit from a single branch to minimize memory usage
	repo, err := git.Clone(store, tree, &git.CloneOptions{
		URL:             repoUrl.String(),
		ReferenceName:   plumbing.NewBranchReferenceName(d.Branch),
		SingleBranch:    true,
		Depth:           1,
		InsecureSkipTLS: !d.HttpsValidate,
	})

	if err != nil {
		result.Message = fmt.Sprintf("Failed to clone the repositoy: %s", err)
		return result
	}

	// Check the contents of the file if configured
	if d.ContentMatch {

		// Open the file for reading
		file, err := tree.Open(d.ContentFile)
		if err != nil {
			result.Message = fmt.Sprintf("Unable to find %s in repository: %s", d.ContentFile, err)
			return result
		}
		defer file.Close()

		// Read the contents of the file
		content, err := io.ReadAll(file)
		if err != nil {
			result.Message = fmt.Sprintf("Failed to read file %s contents: %s", d.ContentFile, err)
			return result
		}

		// Compile the regex statement
		regex, err := regexp.Compile(d.ContentRegex)
		if err != nil {
			result.Message = fmt.Sprintf("Failed to compile regex '%s': %s", d.ContentRegex, err)
			return result
		}

		// Compare the regex with the checked file's contents
		if !regex.Match(content) {
			result.Message = "Matching content not found"
			return result
		}

		// This check has passed, contiune with any other configured checks
	}

	// Check the latest commit's hash
	if d.CommitHashMatch {

		// Get the latest commit
		head, err := repo.Head()
		if err != nil {
			result.Message = fmt.Sprintf("Failed to get where HEAD is point to: %s", err)
			return result
		}

		// Compare the hashes
		if d.CommitHash != head.Hash().String() {
			result.Message = "Commit hash does not match"
			return result
		}

		// This check has passed, contiune with any other configured checks
	}

	// If we reach here, all configured checks have passed
	result.Passed = true
	return result
}

// GetConfig returns the current CheckConfig struct
func (d *Definition) GetConfig() check.Config {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct.
func (d *Definition) SetConfig(c check.Config) {
	d.Config = c
}
