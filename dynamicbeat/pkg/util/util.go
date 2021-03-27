package util

import (
	"bytes"
	"fmt"
	"html/template"
)

func ApplyTemplating(source string, variables map[string]string) (string, error) {
	t := template.New("")
	t, err := t.Parse(source)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %s", err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, variables)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %s", err)
	}

	return buf.String(), nil
}
