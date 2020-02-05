package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of an HTTP check.
type Definition struct {
	ID                 string            // a unique identifier for this check
	Name               string            // a human-readable title for this check
	Host               string            // (required) IP or FQDN of the HTTP server
	Path               string            // (required) path to request - see RFC3986, section 3.3
	HTTPS              bool              // (optional, default false) if HTTPS is to be used
	Port               uint16            // (optional, default 80) TCP port number the HTTP server is listening on
	Method             string            // (optional, default `GET`) HTTP method to use
	Headers            map[string]string // (optional, default empty) name-value pairs of header fields to add/override
	Body               string            // (optional, default empty) the request body
	Verify             bool              // (optional, default false) whether HTTPS certs should be verified
	MatchCode          bool              // (optional, default false) whether the response code must match a defined value for the check to pass
	Code               uint16            // (optional, default 200) the response status code to match
	MatchContent       bool              // (optional, default false) whether the response body must match a defined regex for the check to pass
	ContentRegex       string            // (optional, default `.*`) regex for the response body to match
	SaveMatchedContent bool              // (optional, default false) whether the matched content should be returned in the CheckResult
}

// Run : Execute the check
func Run(chk schema.Check) {
	defer chk.WaitGroup.Done()

	// Set up result
	result := schema.CheckResult{
		Timestamp: time.Now(), // TODO: track how long each check takes
		ID:        chk.ID,
		Name:      chk.Name,
		CheckType: "http",
		Passed:    false,
		Message:   "",
		Details:   nil,
	}

	// Configure HTTP client
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		result.Message = "Could not create CookieJar"
		chk.Output <- result
		return
	}
	client := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: chk.DefinitionList[0]["verify"] == "false", // TODO: document this
			},
		},
	}

	// Save the last returned match string
	var lastMatch *string

	// Make each request in the list
	for _, def := range chk.DefinitionList {
		pass, match, err := request(client, def)

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

	chk.Output <- result
}

func request(client *http.Client, def map[string]string) (bool, *string, error) {
	// Construct URL
	var schema string
	if def["https"] == "true" {
		schema = "https"
	} else {
		schema = "http"
	}
	port := ""
	if portStr, ok := def["port"]; ok {
		port = fmt.Sprintf(":%s", portStr)
	}
	url := fmt.Sprintf("%s://%s%s%s", schema, def["host"], port, def["path"])

	// Construct request
	req, err := http.NewRequest(def["method"], url, strings.NewReader(def["body"]))

	// Add headers
	if gjson.Valid(def["headers"]) {
		headers := gjson.Parse(def["headers"]).Map()
		for k, v := range headers {
			req.Header[k] = []string{v.String()}
		}
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	// Check status code
	if code, ok := def["code"]; ok {
		codeInt, err := strconv.Atoi(code)
		if err != nil {
			return false, nil, fmt.Errorf("Status code must be int: %s", code)
		}
		if resp.StatusCode != codeInt {
			return false, nil, fmt.Errorf("Recieved bad status code: %d", resp.StatusCode)
		}
	}

	// Check body content
	var matchStr string
	if contentMatch, ok := def["content_match"]; ok {
		// Read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, nil, fmt.Errorf("Recieved error when reading response body: %s", err)
		}

		// Check if body matches regex
		regex, err := regexp.Compile(contentMatch)
		if err != nil {
			return false, nil, fmt.Errorf("Error compiling regex string %s : %s", contentMatch, err)
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
func (d Definition) Init(id string, name string, def []byte) error {
	// Set ID and name attributes
	d.ID = id
	d.Name = name

	// Unpack definition json
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}

	// Set nonzero default values
	d.Port = 80
	d.Method = "GET"
	d.Headers = make(map[string]string)
	d.Code = 200
	d.ContentRegex = ".*"

	// Make sure required fields are defined
	missingFields := make([]string, 0)
	if d.Host == "" {
		missingFields = append(missingFields, "Host")
	}

	if d.Path == "" {
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
	return nil
}
