package checksource

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/util"
	"go.uber.org/zap"
)

type Filesystem struct {
	Path  string
	Teams []config.Team
}

func (f *Filesystem) LoadAll() ([]check.Config, error) {
	files, err := os.ReadDir(f.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of directory '%s': %s", f.Path, err)
	}

	// Creating a zero-length slice because we don't know how many of the files
	// in the directory are checks. Also, even if all the files _are_ checks,
	// some of them might not be valid.
	checks := make([]check.Config, 0)
	for _, file := range files {
		if file.IsDir() {
			zap.S().Debugf("skipping directory '%s'", file.Name())
			continue
		}

		filename := file.Name()
		if filepath.Ext(filename) != ".json" {
			zap.S().Debugf("skipping non-JSON file '%s'", file.Name())
			continue
		}

		id := strings.TrimSuffix(filename, ".json")

		// Build the check for each team
		for _, team := range f.Teams {
			fullId := fmt.Sprintf("%s-%s", id, team.Name)
			c, err := f.LoadCheck(fullId)
			if err != nil {
				zap.S().Errorf("skipping check %s due to error when loading: %s", id, err)
			} else {
				checks = append(checks, *c)
			}

		}
	}

	return checks, nil
}

func (f *Filesystem) LoadCheck(id string) (*check.Config, error) {
	// The check ID is made up of the base ID, followed by a '-', then the team
	// name. For example, http-kibana-team01 is a check with a base ID of
	// http-kibana for team01. The base ID is the name of the file the
	// definition is in. For http-kibana-team01, the file would be named
	// http-kibana.json
	s := strings.Split(id, "-")
	teamName := s[len(s)-1]
	baseId := strings.TrimSuffix(id, fmt.Sprintf("-%s", teamName))

	// Find the team struct of the given name
	var team *config.Team
	for _, t := range f.Teams {
		if t.Name == teamName {
			team = &t
			break
		}
	}
	if team == nil {
		return nil, fmt.Errorf("check ID '%s' implies team named '%s', but no team with that name has been configured", id, teamName)
	}

	filepath := fmt.Sprintf("%s%c%s.json", f.Path, os.PathSeparator, baseId)
	body, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read check file '%s': %s", filepath, err)
	}

	// Grab team attribute overrides, if they exist
	overrides := make(map[string]string)
	for k, v := range team.Overrides {
		overrides[k] = v
	}

	// Add attribute for team number if it doesn't exist
	if _, exists := overrides["TeamNum"]; !exists {
		re := regexp.MustCompile(`\S?0*(\d+)$`)
		mat := re.FindStringSubmatch(team.Name)
		overrides["TeamNum"] = mat[len(mat)-1]
	}

	checkFile := struct {
		check.Metadata
		Definition map[string]interface{} `json:"definition"`
		Attributes struct {
			Admin map[string]string `json:"admin"`
			User  map[string]string `json:"user"`
		} `json:"attributes"`
	}{}
	err = json.Unmarshal(body, &checkFile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal check file '%s' from JSON: %s", filepath, err)
	}

	def, err := json.Marshal(checkFile.Definition)
	if err != nil {
		return nil, fmt.Errorf("failed to re-marshal check definition from '%s' to JSON string: %s", filepath, err)
	}

	admin, err := applyOverrides(overrides, checkFile.Attributes.Admin)
	if err != nil {
		return nil, fmt.Errorf("failed to apply overrides to '%s' admin attributes: %s", id, err)
	}
	checkFile.Attributes.Admin = admin
	user, err := applyOverrides(overrides, checkFile.Attributes.User)
	if err != nil {
		return nil, fmt.Errorf("failed to apply overrides to '%s' user attributes: %s", id, err)
	}
	checkFile.Attributes.User = user

	// The ID and group fields are omitted from check definition files
	checkFile.ID = id
	checkFile.Group = teamName

	return &check.Config{
		Metadata:   checkFile.Metadata,
		Definition: def,
		Attributes: check.Attributes{
			Admin: admin,
			User:  user,
		},
	}, nil
}

func applyOverrides(overrides map[string]string, attributes map[string]string) (map[string]string, error) {
	for k, v := range attributes {
		// If the attribute name exists in the overrides map, set its value to
		// whatever's defined in the overrides
		if override, exists := overrides[k]; exists {
			attributes[k] = override
			continue
		}

		// Otherwise, parse the attribute value as a template, using the
		// overrides as keys for the template.
		val, err := util.ApplyTemplating(v, overrides)
		if err != nil {
			return nil, err
		}

		attributes[k] = val
	}

	return attributes, nil
}
