package setup

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
	Inner         http.Client
	Username      string
	Password      string
	Elasticsearch string
	Kibana        string
}

func (c *Client) ReqElasticsearch(method string, path string, body io.Reader) (int, io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s", c.Elasticsearch, path)

	// Build request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to build Elasticsearch request to '%s': %s", path, err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("kbn-xsrf", "true")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Send request
	res, err := c.Inner.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send Elasticsearch request to '%s': %s", path, err)
	}
	return res.StatusCode, res.Body, nil
}

func (c *Client) ReqKibana(method string, path string, body io.Reader) (int, io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s", c.Kibana, path)

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

func (c *Client) Wait() error {
	for {
		_, body, err := c.ReqElasticsearch("GET", "/_cluster/health", nil)
		if err != nil {
			return err
		}

		// Check if response status is "green"
		health := struct {
			Status string `json:"status"`
		}{}
		decoder := json.NewDecoder(body)
		err = decoder.Decode(&health)
		if err != nil {
			return err
		}
		body.Close()
		if health.Status == "green" {
			break
		}

		zap.S().Info("waiting for Elasticsearch to be ready...")
		time.Sleep(5 * time.Second)
	}

	for {
		_, body, err := c.ReqKibana("GET", "/api/status", nil)
		if err != nil {
			return err
		}

		// Check if response status is "green"
		health := struct {
			Status struct {
				Overall struct {
					State string `json:"state"`
				} `json:"overall"`
			} `json:"status"`
		}{}
		decoder := json.NewDecoder(body)
		err = decoder.Decode(&health)
		if err != nil {
			return err
		}
		body.Close()
		if health.Status.Overall.State == "green" {
			break
		}

		zap.S().Info("waiting for Kibana to be ready...")
		time.Sleep(5 * time.Second)
	}

	return nil
}

func CloseAndCheck(code int, body io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer body.Close()
	if code != 200 {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, body)
		if err != nil {
			return fmt.Errorf("got %v response code and couldn't read response body: %s", code, err)
		}
		return fmt.Errorf("response code was %v - response body: %s", code, buf)
	}

	return nil
}

func (c *Client) AddDashboard(data io.Reader) error {
	err := CloseAndCheck(c.ReqKibana("POST", "/api/kibana/dashboards/import?force=true", data))
	if err != nil {
		return err
	}

	return CloseAndCheck(c.ReqKibana("POST", "/s/scorestack/api/kibana/dashboards/import?force=true", data))
}

func (c *Client) AddIndex(name string, data io.Reader) error {
	return CloseAndCheck(c.ReqElasticsearch("PUT", fmt.Sprintf("/_security/role/%s", name), data))
}

func (c *Client) AddRole(name string, data io.Reader) error {
	return CloseAndCheck(c.ReqKibana("PUT", fmt.Sprintf("/api/security/role/%s", name), data))
}

func (c *Client) AddUser(name string, data io.Reader) error {
	url := fmt.Sprintf("/_securty/user/%s", name)

	// Don't try to create the user if they exist already
	code, b, err := c.ReqElasticsearch("GET", url, nil)
	if err != nil {
		return nil
	}
	b.Close()

	if code != 404 {
		return nil
	}

	return CloseAndCheck(c.ReqElasticsearch("PUT", url, data))
}
