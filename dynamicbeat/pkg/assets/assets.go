package assets

import (
	"bytes"
	"embed"
	"io"
	"text/template"
)

const scoreboard = "dashboards/scoreboard.json"
const teamOverview = "dashboards/team-overview.json"
const resultsAdmin = "indices/results-admin.json"
const resultsAll = "indices/results-all.json"
const resultsTeam = "indices/results-team.json"

//go:embed indices dashboards
var f embed.FS

func e(err error) (io.Reader, error) {
	return nil, err
}

func ok(b []byte) (io.Reader, error) {
	return bytes.NewReader(b), nil
}

func read(filename string) (io.Reader, error) {
	data, err := f.ReadFile(filename)
	if err != nil {
		return e(err)
	}

	return ok(data)
}

func Scoreboard() (io.Reader, error) {
	return read(scoreboard)
}

func ResultsAdmin() (io.Reader, error) {
	return read(resultsAdmin)
}

func ResultsAll() (io.Reader, error) {
	return read(resultsAll)
}

func ResultsTeam() (io.Reader, error) {
	return read(resultsTeam)
}

func TeamOverview(name string) (io.Reader, error) {
	data, err := f.ReadFile(teamOverview)
	if err != nil {
		return e(err)
	}

	// Template in the team name
	vars := struct {
		Team string
	}{name}
	tmpl, err := template.New("").Parse(string(data))
	if err != nil {
		return e(err)
	}

	// Apply the template and write to a byte buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		return e(err)
	}

	return ok(buf.Bytes())
}
