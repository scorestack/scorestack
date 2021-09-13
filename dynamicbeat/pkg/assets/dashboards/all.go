package dashboards

import (
	"io"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets"
)

func Scoreboard() io.Reader {
	return assets.Read("dashboards/scoreboard.json")
}

func TeamOverview(name string, checks int) func() io.Reader {
	return func() io.Reader {
		return assets.ReadTeamOverview("dashboards/team-overview.json", name, checks)
	}
}
