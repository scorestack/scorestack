package setup

import (
	"crypto/tls"
	"net/http"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
)

func Run() error {
	c := config.Get()

	// Configure TLS verification based on the Dynamicbeat config setting
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !c.VerifyCerts,
		},
	}

	client := Client{
		Inner:         http.Client{Transport: tr},
		Username:      c.Setup.Username,
		Password:      c.Setup.Password,
		Elasticsearch: c.Elasticsearch,
		Kibana:        c.Setup.Kibana,
	}

	err := client.Wait()
	if err != nil {
		return err
	}

	return nil
}
