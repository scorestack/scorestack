package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of an HTTP check.
type Definition struct {
	ID                   string    // a unique identifier for this check
	Name                 string    // a human-readable title for this check
	Group                string    // the group this check is part of
	ScoreWeight          float64   // the weight that this check has relative to others
	Verify               bool      // (optional, default false) whether HTTPS certs should be validated
	ReportMatchedContent bool      // (optional, default false) whether the matched content should be returned in the CheckResult
	Requests             []Request // a list of requests to make
}

// A Request represents a single HTTP request to make.
type Request struct {
	Host         string            // (required) IP or FQDN of the HTTP server
	Path         string            // (required) path to request - see RFC3986, section 3.3
	HTTPS        bool              // (optional, default false) if HTTPS is to be used
	Port         uint16            // (optional, default 80) TCP port number the HTTP server is listening on
	Method       string            // (optional, default `GET`) HTTP method to use
	Headers      map[string]string // (optional, default empty) name-value pairs of header fields to add/override
	Body         string            // (optional, default empty) the request body
	MatchCode    bool              // (optional, default false) whether the response code must match a defined value for the check to pass
	Code         int               // (optional, default 200) the response status code to match
	MatchContent bool              // (optional, default false) whether the response body must match a defined regex for the check to pass
	ContentRegex string            // (optional, default `.*`) regex for the response body to match
	StoreValue   bool              // (optional, default false) whether the matched content should be saved for use in a later request
}

// Run a single instance of the check.
func (d *Definition) Run(ctx context.Context) schema.CheckResult {

	// Set up result
	result := schema.CheckResult{
		Timestamp:   time.Now(),
		ID:          d.ID,
		Name:        d.Name,
		Group:       d.Group,
		ScoreWeight: d.ScoreWeight,
		CheckType:   "http",
	}

	// Configure HTTP client
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		result.Message = "Could not create CookieJar"
		return result
	}
	client := &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !d.Verify,
			},
		},
	}

	// Save match strings
	var lastMatch *string
	var storedValue *string

	type storedValTempl struct {
		SavedValue string
	}

	// Make each request in the list
	for _, r := range d.Requests {
		// Check to see if the StoredValue needs to be templated in
		if storedValue != nil {
			// TODO: refactor this out into a function that keeps the "happy line"
			// Re-encode definition to JSON string
			def, err := json.Marshal(r)
			if err != nil {
				logp.Info("Error encoding HTTP definition as JSON for StoredValue templating: %s", err)
			} else {
				attrs := storedValTempl{
					SavedValue: *storedValue,
				}
				templ := template.Must(template.New("http-storedvalue").Parse(string(def)))
				var buf bytes.Buffer
				err := templ.Execute(&buf, attrs)
				if err != nil {
					logp.Info("Error templating HTTP definition for StoredValue templating: %s", err)
				} else {
					newReq := Request{}
					err := json.Unmarshal(buf.Bytes(), &newReq)
					if err != nil {
						logp.Info("Error decoding StoredValue-templated HTTP definition: %s", err)
					} else {
						r = newReq
					}
				}
			}
		}

		pass, match, err := request(ctx, client, r)

		// Process request results
		result.Passed = pass
		if err != nil {
			result.Message = fmt.Sprintf("%s", err)
		}
		if match != nil {
			lastMatch = match
			if r.StoreValue {
				storedValue = match
			}
		}

		// If this request failed, don't continue on to the next request
		if !pass {
			break
		}
	}

	details := make(map[string]string)
	if d.ReportMatchedContent && lastMatch != nil {
		details["matched_content"] = *lastMatch
	}
	result.Details = details

	result.Passed = true
	return result
}

func request(ctx context.Context, client *http.Client, r Request) (bool, *string, error) {
	// Construct URL
	var schema string
	if r.HTTPS {
		schema = "https"
	} else {
		schema = "http"
	}
	url := fmt.Sprintf("%s://%s:%d%s", schema, r.Host, r.Port, r.Path)

	// Construct request
	req, err := http.NewRequestWithContext(ctx, r.Method, url, strings.NewReader(r.Body))

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
		matches := regex.FindSubmatch(body)
		matchStr = fmt.Sprintf("%s", matches[len(matches)-1])
	}

	// If we've reached this point, then the check succeeded
	return true, &matchStr, nil
}

// Init the check using a known ID and name. The rest of the check fields will
// be filled in by parsing a JSON string representing the check definition.
func (d *Definition) Init(id string, name string, group string, scoreWeight float64, def []byte) error {
	// Unpack definition json
	err := json.Unmarshal(def, &d)
	if err != nil {
		return err
	}
	// TODO: set verify value

	// Set generic attributes
	d.ID = id
	d.Name = name
	d.Group = group
	d.ScoreWeight = scoreWeight

	// Finish initializing each request
	for i := range d.Requests {
		// Set nonzero default values
		if d.Requests[i].Port == 0 {
			d.Requests[i].Port = 80
		}

		if d.Requests[i].Method == "" {
			d.Requests[i].Method = "GET"
		}

		if d.Requests[i].Headers == nil {
			d.Requests[i].Headers = make(map[string]string)
		}

		if d.Requests[i].Code == 0 {
			d.Requests[i].Code = 200
		}

		if d.Requests[i].ContentRegex == "" {
			d.Requests[i].ContentRegex = ".*"
		}

		// Make sure required fields are defined
		missingFields := make([]string, 0)
		if d.Requests[i].Host == "" {
			missingFields = append(missingFields, "Host")
		}

		if d.Requests[i].Path == "" {
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
