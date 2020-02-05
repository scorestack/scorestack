package http

import (
	"crypto/tls"
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
