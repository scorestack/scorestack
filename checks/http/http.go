package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"sync"
	"time"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of an HTTP check.
type Definition struct {
	ID       string    // a unique identifier for this check
	Name     string    // a human-readable title for this check
	Verify   bool      // (optional, default false) whether HTTPS certs should be validated
	Requests []Request // a list of requests to make
}

// A Request represents a single HTTP request to make.
type Request struct {
	Host               string            // (required) IP or FQDN of the HTTP server
	Path               string            // (required) path to request - see RFC3986, section 3.3
	HTTPS              bool              // (optional, default false) if HTTPS is to be used
	Port               uint16            // (optional, default 80) TCP port number the HTTP server is listening on
	Method             string            // (optional, default `GET`) HTTP method to use
	Headers            map[string]string // (optional, default empty) name-value pairs of header fields to add/override
	Body               string            // (optional, default empty) the request body
	MatchCode          bool              // (optional, default false) whether the response code must match a defined value for the check to pass
	Code               int               // (optional, default 200) the response status code to match
	MatchContent       bool              // (optional, default false) whether the response body must match a defined regex for the check to pass
	ContentRegex       string            // (optional, default `.*`) regex for the response body to match
	SaveMatchedContent bool              // (optional, default false) whether the matched content should be returned in the CheckResult
}

// Run a single instance of the check.
func (d *Definition) Run(wg *sync.WaitGroup, out chan<- schema.CheckResult) {
	defer wg.Done()

	// Set up result
	result := schema.CheckResult{
		Timestamp: time.Now(),
		ID:        d.ID,
		Name:      d.Name,
		CheckType: "http",
	}

	// Configure HTTP client
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		result.Message = "Could not create CookieJar"
		out <- result
		return
	}
	client := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !d.Verify,
			},
		},
	}

	// Save the last returned match string
	var lastMatch *string

	// Make each request in the list
	for _, r := range d.Requests {
		pass, match, err := request(client, r)

		// Process request results
		result.Passed = pass
		if err != nil {
			result.Message = fmt.Sprintf("%s", err)
		}
		if match != nil {
			lastMatch = match
		}

		// If this request failed, don't continue on to the next request
		if !pass {
			break
		}
	}

	details := make(map[string]string)
	if lastMatch != nil {
		details["matched_content"] = *lastMatch
	}
	result.Details = details

	out <- result
}

func request(client *http.Client, r Request) (bool, *string, error) {
	// Construct URL
	var schema string
	if r.HTTPS {
		schema = "https"
	} else {
		schema = "http"
	}
	url := fmt.Sprintf("%s://%s:%d%s", schema, r.Host, r.Port, r.Path)

	// Construct request
	req, err := http.NewRequest(r.Method, url, strings.NewReader(r.Body))

	// Add headers
	for k, v := range r.Headers {
		req.Header[k] = []string{v}
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	// Check status code
	if r.MatchCode && resp.StatusCode != r.Code {
		return false, nil, fmt.Errorf("Recieved bad status code: %d", resp.StatusCode)
	}

	// Check body content
	var matchStr string
	if r.MatchContent {
		// Read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, nil, fmt.Errorf("Recieved error when reading response body: %s", err)
		}

		// Check if body matches regex
		regex, err := regexp.Compile(r.ContentRegex)
		if err != nil {
			return false, nil, fmt.Errorf("Error compiling regex string %s : %s", r.ContentRegex, err)
		}
		if !regex.Match(body) {
			return false, nil, fmt.Errorf("recieved bad response body")
		}
		matchStr = fmt.Sprintf("%s", regex.Find(body))
	}

	// If we've reached this point, then the check succeeded
	return true, &matchStr, nil
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, def []byte) error {
	// Set ID and name attributes
	d.ID = id
	d.Name = name

	// Unpack definition json
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}
	// TODO: set verify value

	// Finish initializing each request
	for _, r := range d.Requests {
		// Set nonzero default values
		r.Port = 80
		r.Method = "GET"
		r.Headers = make(map[string]string)
		r.Code = 200
		r.ContentRegex = ".*"

		// Make sure required fields are defined
		missingFields := make([]string, 0)
		if r.Host == "" {
			missingFields = append(missingFields, "Host")
		}

		if r.Path == "" {
			missingFields = append(missingFields, "Path")
		}

		// Error on only the first missing field, if there are any
		if len(missingFields) > 0 {
			return schema.ValidationError{
				ID:    d.ID,
				Type:  "http",
				Field: missingFields[0],
			}
		}
	}

	return nil
}
