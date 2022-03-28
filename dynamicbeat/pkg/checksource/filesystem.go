package checksource

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/config"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/models"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/util"
	"go.uber.org/zap"
)

type Filesystem struct {
	Path  string
	Teams []config.Team
}

func (f *Filesystem) Error(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s directory '%s': %s", msg, f.Path, err)
	}

	return fmt.Errorf("%s directory '%s'", msg, f.Path)
}

func (f *Filesystem) LoadAll() ([]models.CheckConfig, error) {
	fInfo, err := os.Stat(f.Path)
	if err != nil {
		return nil, f.Error("failed to stat", err)
	}

	if fInfo.IsDir() {
		zap.S().Infof("searching for checks within directory %s", f.Path)
		configs := make([]models.CheckConfig, 0)

		matches, err := filepath.Glob(filepath.Join(f.Path, "**", "*.toml"))
		if err != nil {
			return nil, f.Error("failed to glob for TOML files in", err)
		}
		for _, match := range matches {
			fileConfigs, err := f.LoadFile(match)
			if err != nil {
				msg := fmt.Sprintf("skipping check file '%s' in", match)
				zap.S().Error(f.Error(msg, err))
				continue
			}

			configs = append(configs, fileConfigs...)
		}

		return configs, nil
	}

	return f.LoadFile(f.Path)
}

// LoadFile: load a check file and create a check for each team
func (f *Filesystem) LoadFile(path string) ([]models.CheckConfig, error) {
	zap.S().Infof("loading check from %s", path)

	var source TomlConfig
	_, err := toml.DecodeFile(path, &source)
	if err != nil {
		return nil, fmt.Errorf("failed to decode check file '%s': %s", path, err)
	}
	// TODO: validate Kind

	// Default to 1 point if points are unspecified
	if source.points == 0 {
		source.points = 1
	}

	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	// Create the check for each team
	configs := make([]models.CheckConfig, 0)
	for _, team := range f.Teams {
		config, err := source.TeamConfig(baseName, team)
		if err != nil {
			zap.S().Errorf("failed to create config from file '%s' for team '%s'", path, team.Name)
			return nil, err
		}

		configs = append(configs, config)
	}

	return configs, nil
}

type TomlConfig struct {
	display_name string
	description  string
	kind         string
	points       uint64
	definition   map[string]interface{}
	attribute    []struct {
		key          string
		value        string
		display_name string
		permissions  models.Permission
		display_as   models.ViewType
		description  string
	}
}

func (t *TomlConfig) TeamConfig(filename string, team config.Team) (models.CheckConfig, error) {
	config := models.CheckConfig{
		CheckMetadata: models.CheckMetadata{
			CheckId:     filename,
			Group:       team.Name,
			DisplayName: t.display_name,
			Description: t.description,
			Kind:        t.kind,
			Points:      t.points,
		},
		Attributes: make([]models.Attribute, 0),
	}

	// Grab any defined team overrides
	overrides := make(map[string]string)
	for k, v := range team.Overrides {
		overrides[k] = v
	}

	// Add overrides for team name and number if they don't exist
	if _, exists := overrides["TeamNum"]; !exists {
		re := regexp.MustCompile(`\S?0*(\d+)$`)
		mat := re.FindStringSubmatch(team.Name)
		overrides["TeamNum"] = mat[len(mat)-1]
	}
	if _, exists := overrides["TeamName"]; !exists {
		overrides["TeamName"] = team.Name
	}

	// Create each attribute
	for _, a := range t.attribute {
		meta := models.AttributeMetadata{
			CheckId:     config.CheckId,
			Group:       config.Group,
			Key:         a.key,
			Permissions: a.permissions,
		}

		attribute := models.Attribute{
			AttributeMetadata: meta,
			DisplayName:       a.display_name,
			Description:       a.description,
			DisplayAs:         a.display_as,
			Values:            make([]models.AttributeValue, 0),
		}

		v := models.AttributeValue{
			AttributeMetadata: meta,
			Created:           time.Now(),
		}

		// Apply overrides to the attribute value
		value, err := util.ApplyTemplating(a.value, overrides)
		if err != nil {
			zap.S().Error(attribute.Error("skipping overrides", err))

			// Fall back to the original value
			value = a.value
		}
		v.Value = value

		attribute.Values = append(attribute.Values, v)
		config.Attributes = append(config.Attributes, attribute)
	}

	return config, nil
}
