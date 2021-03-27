package esclient

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go.uber.org/zap"
)

func (c *Client) CloseAndCheck(res *esapi.Response) error {
	defer res.Body.Close()
	if res.IsError() {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, res.Body)
		if err != nil {
			return fmt.Errorf("got %v response code and couldn't read response body: %s", res.StatusCode, err)
		}
		return fmt.Errorf("response code was %v - response body: %s", res.StatusCode, buf)
	}

	return nil
}

func (c *Client) Wait() error {
	first := true
	for {
		// If we haven't been through this loop yet, sleep for 5 seconds
		if !first {
			zap.S().Info("waiting for Elasticsearch to be ready...")
			time.Sleep(5 * time.Second)
		}
		first = false

		// We could use WithWaitForStatus to block until the status is green,
		// but then there wouldn't be output to the user that we were still
		// waiting. Waiting for green status takes a while, and if we didn't
		// periodically update the user that we're still waiting, they might
		// get concerned that the program isn't working.
		res, err := c.Cluster.Health()
		if err != nil || res.IsError() {
			continue
		}

		// Check if response status is "green"
		health := struct {
			Status string `json:"status"`
		}{}
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&health)
		if err != nil {
			continue
		}
		res.Body.Close()
		if health.Status == "green" {
			break
		}
	}

	return nil
}

func (c *Client) AddIndex(name string, body io.Reader) error {
	res, err := c.Indices.Exists([]string{name})
	if err != nil {
		return fmt.Errorf("failed to check if index '%s' exists: %s", name, err)
	}
	defer res.Body.Close()

	// Don't create the index if it already exists
	if res.StatusCode == 200 {
		zap.S().Debugf("index '%s' already exists, skipping...", name)
		return nil
	}

	res, err = c.Indices.Create(name, c.Indices.Create.WithBody(body))
	if err != nil {
		return fmt.Errorf("failed to create index '%s': %s", name, err)
	}

	return c.CloseAndCheck(res)
}

func (c *Client) AddUser(name string, body io.Reader) error {
	res, err := c.Security.GetUser(c.Security.GetUser.WithUsername(name))
	if err != nil {
		return fmt.Errorf("failed to check if user '%s' exists: %s", name, err)
	}
	defer res.Body.Close()

	// Don't create the user if they already exist
	if res.StatusCode == 200 {
		zap.S().Debugf("user '%s' already exists, skipping...", name)
		return nil
	}

	res, err = c.Security.PutUser(name, body)
	if err != nil {
		return fmt.Errorf("failed to create user '%s': %s", name, err)
	}

	return c.CloseAndCheck(res)
}
