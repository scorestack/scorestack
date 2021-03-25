package setup

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/dashboards"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/indices"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/roles"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/spaces"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"go.uber.org/zap"
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

	zap.S().Info("checking if Elasticsearch and Kibana are up")
	err := client.Wait()
	if err != nil {
		return err
	}

	err = client.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Initialize() error {
	// Add Dynamicbeat role and user
	err := c.AddRole("dynamicbeat", roles.Dynamicbeat())
	if err != nil {
		return err
	}
	err = c.AddUser("dynamicbeat", roles.Dynamicbeat())
	if err != nil {
		return err
	}

	zap.S().Info("adding Scorestack space")
	err = CloseAndCheck(c.ReqKibana("PUT", "/api/spaces/space/scorestack", spaces.Scorestack()))
	if err != nil {
		return err
	}

	zap.S().Info("enabling dark theme")
	valTrue := strings.NewReader(`{"value":"true"}`)
	err = CloseAndCheck(c.ReqKibana("POST", "/api/kibana/settings/theme:darkMode", valTrue))
	if err != nil {
		return err
	}
	valTrue = strings.NewReader(`{"value":"true"}`)
	err = CloseAndCheck(c.ReqKibana("POST", "/s/scorestack/api/kibana/settings/theme:darkMode", valTrue))
	if err != nil {
		return err
	}

	// Add base role for common permissions
	err = c.AddRole("common", roles.Common())
	if err != nil {
		return err
	}

	// Add spectator role
	err = c.AddRole("spectator", roles.Spectator())
	if err != nil {
		return err
	}

	// Add admin roles
	err = c.AddRole("attribute-admin", roles.AttributeAdmin())
	if err != nil {
		return err
	}
	err = c.AddRole("check-admin", roles.AttributeAdmin())
	if err != nil {
		return err
	}

	// Add Scoreboard dashboard
	err = c.AddDashboard(dashboards.Scoreboard)
	if err != nil {
		return err
	}

	// Add default index template
	zap.S().Info("adding default index template")
	idx := strings.NewReader(`{"index_patterns":["check*","attrib_*","results*"],"settings":{"number_of_replicas":"0"}}`)
	err = CloseAndCheck(c.ReqElasticsearch("PUT", "/_template/default", idx))
	if err != nil {
		return err
	}

	// Create results indices
	err = c.AddIndex("results-admin", indices.ResultsAdmin())
	if err != nil {
		return err
	}
	err = c.AddIndex("results-all", indices.ResultsAll())
	if err != nil {
		return err
	}

	return nil
}
