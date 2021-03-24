package esclient

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

func New(host string, username string, password string, verify bool) (*elasticsearch.Client, error) {
	clientConfig := elasticsearch.Config{
		Addresses: []string{host},
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			DialContext:         (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !verify,
			},
		},
	}
	es, err := elasticsearch.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %s", err)
	}

	return es, nil
}
