package setup

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/dashboards"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/indices"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/roles"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/spaces"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/users"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checksource"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/esclient"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/kibclient"
	"go.uber.org/zap"
)

func Run() error {
	c := config.Get()

	err := kibSetup(c.Setup.Kibana, c.Setup.Username, c.Setup.Password, c.VerifyCerts)
	if err != nil {
		return err
	}

	es, err := esclient.New(c.Elasticsearch, c.Setup.Username, c.Setup.Password, c.VerifyCerts)
	if err != nil {
		return err
	}

	err = esSetup(es)
	if err != nil {
		return err
	}

	f := &checksource.Filesystem{
		Path:  c.Setup.CheckFolder,
		Teams: c.Teams,
	}
	defs, err := f.LoadAll()
	if err != nil {
		return err
	}

	for _, def := range defs {
		fmt.Printf("%s: %#v", def.ID, def)
	}

	return nil
}

func kibSetup(host string, user string, pass string, verify bool) error {
	// Configure TLS verification based on the Dynamicbeat config setting
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !verify,
		},
	}

	c := kibclient.Client{
		Inner:    http.Client{Transport: tr, Timeout: 5 * time.Second},
		Username: user,
		Password: pass,
		Host:     host,
	}

	zap.S().Info("checking if Kibana is up")
	err := c.Wait()
	if err != nil {
		return err
	}

	// Add Dynamicbeat role
	err = c.AddRole("dynamicbeat", roles.Dynamicbeat())
	if err != nil {
		return err
	}

	// Add Scorestack space
	err = c.AddSpace("scorestack", spaces.Scorestack)
	if err != nil {
		return err
	}

	zap.S().Info("enabling dark theme")
	valTrue := strings.NewReader(`{"value":"true"}`)
	err = c.CheckedReq("POST", "/api/kibana/settings/theme:darkMode", valTrue)
	if err != nil {
		return err
	}
	valTrue = strings.NewReader(`{"value":"true"}`)
	err = c.CheckedReq("POST", "/s/scorestack/api/kibana/settings/theme:darkMode", valTrue)
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

	return nil
}

func esSetup(c *esclient.Client) error {
	zap.S().Info("checking if Elasticsaerch is up")
	err := c.Wait()
	if err != nil {
		return err
	}

	err = c.AddUser("dynamicbeat", users.Dynamicbeat())
	if err != nil {
		return err
	}

	// Add default index template
	zap.S().Info("adding default index template")
	idx := strings.NewReader(`{"index_patterns":["check*","attrib_*","results*"],"settings":{"number_of_replicas":"0"}}`)
	res, err := c.Indices.PutTemplate("default", idx)
	if err != nil {
		return err
	}
	err = c.CloseAndCheck(res)
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
