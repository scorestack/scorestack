package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/common"
)

// Run : Execute the check
func Run(chk common.Check) {
	defer chk.WaitGroup.Done()

	startTime := time.Now()

	// Set up result
	result := common.CheckResult{
		Timestamp: startTime,
		ID:        chk.ID,
		Name:      chk.Name,
		CheckType: "http",
		Passed:    false,
		Message:   "",
		Details:   chk.Definition,
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
				InsecureSkipVerify: chk.Definition["verify"] == "false",
			},
		},
	}

	// Construct URL
	url := make([]string, 0)
	if chk.Definition["https"] == "true" {
		url = append(url, "https://")
	} else {
		url = append(url, "http://")
	}
	url = append(url, chk.Definition["host"])
	if _, ok := chk.Definition["port"]; ok {
		url = append(url, ":")
		url = append(url, chk.Definition["port"])
	}
	url = append(url, chk.Definition["path"])
	urlString := strings.Join(url, "")

	details := make(map[string]string)
	switch chk.Definition["method"] {
	case "GET":
		resp, err := client.Get(urlString)
		if err != nil {
			result.Message = fmt.Sprintf("Error with GET: %s", err)
			chk.Output <- result
			return
		}
		defer resp.Body.Close()
		details["code"] = string(resp.StatusCode)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		details["body"] = string(bodyBytes)
	}

	result.Details = details

	chk.Output <- result
}
