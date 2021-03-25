package dashboards

import (
	"io"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets"
)

func Scoreboard() io.Reader {
	return assets.Read("dashboards/scoreboard.json")
}

func TeamOverview(name string) io.Reader {
	return assets.ReadTeam("dashboards/team-overview.json", name)
}
