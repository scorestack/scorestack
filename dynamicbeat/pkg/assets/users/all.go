package users

import (
	"io"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/assets"
)

func Dynamicbeat() io.Reader {
	return assets.Read("users/dynamicbeat.json")
}

func Team(name string) io.Reader {
	return assets.ReadTeam("users/team.json", name)
}
