package checksource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"go.uber.org/zap"
)

type Filesystem struct {
	Path  string
	Teams []string
}

func (f *Filesystem) LoadAll() ([]check.Config, error) {
	files, err := ioutil.ReadDir(f.Path)
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

		c, err := f.LoadCheck(id)
		if err != nil {
			zap.S().Errorf("skipping check %s due to error when loading: %s", id, err)
		}

		// Build the check for each team
		for _, team := range f.Teams {
			teamCheck := c
			teamCheck.ID = fmt.Sprintf("%s-%s", teamCheck.ID, team)
			teamCheck.Group = team

			// Add attribute for team number
			re := regexp.MustCompile(`\S?0*(\d+)$`)
			mat := re.FindStringSubmatch(team)
			teamNum := mat[len(mat)-1]
			teamCheck.Attribs["TeamNum"] = teamNum

			checks = append(checks, *c)
		}
	}

	return checks, nil
}

func (f *Filesystem) LoadCheck(id string) (*check.Config, error) {

	filepath := fmt.Sprintf("%s/%s.json", f.Path, id)
	body, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read check file '%s': %s", filepath, err)
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

	// Marge attributes
	attr := make(map[string]string)
	for k, v := range checkFile.Attributes.Admin {
		attr[k] = v
	}
	for k, v := range checkFile.Attributes.User {
		attr[k] = v
	}

	return &check.Config{
		Metadata:   checkFile.Metadata,
		Definition: def,
		Attribs:    attr,
	}, nil
}
