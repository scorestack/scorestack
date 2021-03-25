package setup

import (
	"encoding/json"
	"fmt"
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

func (c *Client) Wait() error {
	for {
		ready, err := c.esIsReady()
		if err != nil {
			return err
		}

		if ready {
			break
		}

		zap.S().Info("waiting for Elasticsearch to be ready...")
		time.Sleep(5 * time.Second)
	}

	for {
		ready, err := c.kibIsReady()
		if err != nil {
			return err
		}

		if ready {
			break
		}

		zap.S().Info("waiting for Kibana to be ready...")
		time.Sleep(5 * time.Second)
	}

	return nil
}

type clusterHealth struct {
	Status string `json:"status"`
}

func (c *Client) esIsReady() (bool, error) {
	url := fmt.Sprintf("%s/_cluster/health", c.Elasticsearch)

	// Build request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to build Elasticsearch health request: %s", err)
	}
	req.SetBasicAuth(c.Username, c.Password)

	// Send request
	res, err := c.Inner.Do(req)
	if err != nil {
		return false, nil
	}
	defer res.Body.Close()

	// Check if response status is "green"
	var health clusterHealth
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&health)
	if health.Status == "green" {
		return true, nil
	} else {
		return false, nil
	}
}

type kibanaStatus struct {
	Status struct {
		Overall struct {
			State string `json:"state"`
		} `json:"overall"`
	} `json:"status"`
}

func (c *Client) kibIsReady() (bool, error) {
	url := fmt.Sprintf("%s/api/status", c.Kibana)

	// Build request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to build Kibana status request: %s", err)
	}
	req.SetBasicAuth(c.Username, c.Password)

	// Send request
	res, err := c.Inner.Do(req)
	if err != nil {
		return false, nil
	}
	defer res.Body.Close()

	// Check if response status is "green"
	var health kibanaStatus
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&health)
	if health.Status.Overall.State == "green" {
		return true, nil
	} else {
		return false, nil
	}
}
