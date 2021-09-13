package setup

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/dashboards"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/roles"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets/spaces"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/kibclient"
	"go.uber.org/zap"
)

func Kibana(host string, user string, pass string, verify bool, teams []config.Team) error {
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

	for _, team := range teams {
		err = c.AddRole(team.Name, roles.Team(team.Name))
		if err != nil {
			zap.S().Errorf("failed to add role for %s: %s", team.Name, err)
		}

		// TODO: don't hardcode the number of rows in the table
		err = c.AddDashboard(dashboards.TeamOverview(team.Name, 20))
		if err != nil {
			zap.S().Errorf("failed to add team overview dashboard for %s: %s", team.Name, err)
		}
	}

	return nil
}
