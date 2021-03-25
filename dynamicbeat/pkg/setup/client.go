package setup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
