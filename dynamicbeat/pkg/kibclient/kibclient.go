package kibclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	Inner    http.Client
	Username string
	Password string
	Host     string
}

func (c *Client) Req(method string, path string, body io.Reader) (int, io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s", c.Host, path)

	// Build request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to build Kibana request to '%s': %s", path, err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("kbn-xsrf", "true")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Send request
	res, err := c.Inner.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send Kibana request to '%s': %s", path, err)
	}
	return res.StatusCode, res.Body, nil
}

func (c *Client) CheckedReq(method string, path string, body io.Reader) error {
	return CloseAndCheck(c.Req(method, path, body))
}

func (c *Client) Wait() error {
	first := true
	for {
		// If we haven't been through this loop yet, sleep for 5 seconds
		if !first {
			zap.S().Info("waiting for Kibana to be ready...")
			time.Sleep(5 * time.Second)
		}
		first = false

		_, body, err := c.Req("GET", "/api/status", nil)
		if err != nil {
			continue
		}

		// Check if response status is "green"
		health := struct {
			Status struct {
				Overall struct {
					State string `json:"state"`
					Level string `json:"level"`
				} `json:"overall"`
			} `json:"status"`
		}{}
		decoder := json.NewDecoder(body)
		err = decoder.Decode(&health)
		if err != nil {
			continue
		}
		body.Close()
		if health.Status.Overall.State == "green" || health.Status.Overall.Level == "available" {
			break
		}
	}

	return nil
}

func CloseAndCheck(code int, body io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer body.Close()
	if code != 200 && code != 204 {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, body)
		if err != nil {
			return fmt.Errorf("got %v response code and couldn't read response body: %s", code, err)
		}
		return fmt.Errorf("response code was %v - response body: %s", code, buf)
	}

	return nil
}

func (c *Client) AddDashboard(data func() io.Reader) error {
	zap.S().Info("adding dashboards")
	err := CloseAndCheck(c.Req("POST", "/api/kibana/dashboards/import?force=true", data()))
	if err != nil {
		return err
	}

	return CloseAndCheck(c.Req("POST", "/s/scorestack/api/kibana/dashboards/import?force=true", data()))
}

func (c *Client) AddRole(name string, data io.Reader) error {
	zap.S().Infof("adding role: %s", name)
	return CloseAndCheck(c.Req("PUT", fmt.Sprintf("/api/security/role/%s", name), data))
}

func (c *Client) AddSpace(name string, data func() io.Reader) error {
	// Try to update the space if it already exists
	code, b, err := c.Req("PUT", fmt.Sprintf("/api/spaces/space/%s", name), data())
	if code == 404 {
		// If the space doesn't exist, create it
		zap.S().Infof("adding Kibana space: %s", name)
		return CloseAndCheck(c.Req("POST", "/api/spaces/space", data()))
	}

	zap.S().Debugf("Kibana space '%s' already exists, skipping...", name)
	return CloseAndCheck(code, b, err)
}
