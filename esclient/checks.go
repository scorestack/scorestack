package esclient

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/schema"
)

// UpdateCheckDefs will re-read all check definitions from a single index and
// load the related attributes for each check.
func UpdateCheckDefs(c *elasticsearch.Client, i string) ([]schema.CheckDef, error) {
	results := make([]schema.CheckDef, 0)

	// Get list of checks
	checks, err := GetAllDocuments(c, i)
	if err != nil {
		return nil, err
	}

	// Iterate over each check
	for _, check := range checks {
		// Decode check definition
		checkMap := make(map[string]interface{})
		// We can assume that the JSON can be unmarshalled, because the JSON
		// was created with json.Marshal()
		_ = json.Unmarshal(check, &checkMap)
		result := schema.CheckDef{
			ID:      checkMap["id"].(string),
			Name:    checkMap["name"].(string),
			Type:    checkMap["type"].(string),
			Attribs: make(map[string]string),
		}

		// Get any template variables for the check
		for _, perm := range []string{"admin", "user"} {
			// Generate attribute index name
			idx := strings.Join([]string{"attrib_", perm, "_", result.ID}, "")

			// Get attribute document
			docs, err := GetAllDocuments(c, idx)
			if err != nil {
				return results, err
			}

			// Decode attribute document
			err = json.Unmarshal(docs[0], &result.Attribs)
			if err != nil {
				return nil, fmt.Errorf("Failed to decode attribute document %s: %s", idx, err)
			}
		}

		results = append(results, result)
	}

	return results, nil
}
