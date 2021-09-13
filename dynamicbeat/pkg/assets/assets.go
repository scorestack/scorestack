package assets

import (
	"bytes"
	"embed"
	"io"
	"text/template"

	"go.uber.org/zap"
)

//go:embed *
var f embed.FS

func Read(filename string) io.Reader {
	data, err := f.ReadFile(filename)
	if err != nil {
		zap.S().Panicf("failed to read embedded asset %s: %s", filename, err)
	}

	return bytes.NewReader(data)
}

func ReadTeam(filename string, name string) io.Reader {
	data, err := f.ReadFile(filename)
	if err != nil {
		zap.S().Panicf("failed to read embedded asset %s: %s", filename, err)
	}

	// Template in the team name
	vars := struct {
		Team string
	}{name}
	tmpl, err := template.New("").Parse(string(data))
	if err != nil {
		zap.S().Panicf("failed to read asset %s as template: %s", filename, err)
	}

	// Apply the template and write to a byte buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		zap.S().Panicf("failed to template team %s into asset %s: %s", name, filename, err)
	}

	return bytes.NewReader(buf.Bytes())
}

func ReadTeamOverview(filename string, name string, checks int) io.Reader {
	data, err := f.ReadFile(filename)
	if err != nil {
		zap.S().Panicf("failed to read embedded asset %s: %s", filename, err)
	}

	// Template in the team name and the number of checks per team
	vars := struct {
		Team   string
		Checks int
	}{name, checks}
	tmpl, err := template.New("").Parse(string(data))
	if err != nil {
		zap.S().Panicf("failed to read asset %s as template: %s", filename, err)
	}

	// Apply the template and write to a byte buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		zap.S().Panicf("failed to template team %s into asset %s: %s", name, filename, err)
	}

	return bytes.NewReader(buf.Bytes())
}
