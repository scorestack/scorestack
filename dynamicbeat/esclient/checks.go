package esclient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"

	"github.com/scorestack/scorestack/dynamicbeat/checks/schema"
)

// UpdateCheckDefs will re-read all check definitions from a single index and
// load the related attributes for each check.
func UpdateCheckDefs(c *elasticsearch.Client, i string) ([]schema.CheckConfig, error) {
	// Track how long it takes to update check definitions
	start := time.Now()

	results := make([]schema.CheckConfig, 0)

	// Get list of checks
	checks, err := GetAllDocuments(c, i)
	if err != nil {
		return nil, err
	}

	// Get list of attributes
	attribDocs, err := GetAllDocuments(c, "attrib_*")
	if err != nil {
		return nil, err
	}

	// Organize attributes by check ID
	attribs := make(map[string][]Document)
	for _, doc := range attribDocs {
		if _, ok := attribs[doc.ID]; !ok {
			attribs[doc.ID] = make([]Document, 0)
		}
		attribs[doc.ID] = append(attribs[doc.ID], doc)
	}

	// Iterate over each check
	for _, doc := range checks {
		// Encode the definition as a JSON string
		def, err := json.Marshal(doc.Source["definition"])
		if err != nil {
			return nil, fmt.Errorf("Error encoding definition for %s to JSON string: %s", doc.ID, err)
		}

		// Unpack check definition into CheckConfig struct
		result := schema.CheckConfig{
			ID:          doc.Source["id"].(string),
			Name:        doc.Source["name"].(string),
			Type:        doc.Source["type"].(string),
			Group:       doc.Source["group"].(string),
			ScoreWeight: doc.Source["score_weight"].(float64),
			Definition:  def,
			Attribs:     make(map[string]string),
		}

		// Add any template variables to the check
		if val, ok := attribs[result.ID]; ok {
			// Decode each attribute in each document
			for _, doc := range val {
				for k, v := range doc.Source {
					// Decode the value of the attribute
					result.Attribs[k] = v.(string)
				}
			}
		}

		// Add the SavedValue attribute in case the check uses it
		result.Attribs["SavedValue"] = "{{.SavedValue}}"

		results = append(results, result)
	}

	logp.Info("Updated check definitions in %.2f seconds", time.Since(start).Seconds())
	return results, nil
}
