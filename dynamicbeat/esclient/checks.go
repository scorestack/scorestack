package esclient

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch"

	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// UpdateCheckDefs will re-read all check definitions from a single index and
// load the related attributes for each check.
func UpdateCheckDefs(c *elasticsearch.Client, i string) ([]schema.CheckConfig, error) {
	results := make([]schema.CheckConfig, 0)

	// Get list of checks
	checks, err := GetAllDocuments(c, i)
	if err != nil {
		return nil, err
	}

	// Get list of attributes
	attribDocs, err := GetAllDocuments(c, "attrib_")
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
	for _, check := range checks {
		// Decode check definition
		checkMap := make(map[string]interface{})
		err = json.Unmarshal(check.Source, &checkMap)
		if err != nil {
			return nil, fmt.Errorf("Error decoding JSON string for definition of %s: %s", check.ID, err)
		}

		// Re-encode definition to JSON string
		def, err := json.Marshal(checkMap["definition"])
		if err != nil {
			return nil, fmt.Errorf("Error encoding definition as JSON: %s", err)
		}

		result := schema.CheckConfig{
			ID:          checkMap["id"].(string),
			Name:        checkMap["name"].(string),
			Type:        checkMap["type"].(string),
			Group:       checkMap["group"].(string),
			ScoreWeight: checkMap["score_weight"].(float64),
			Definition:  def,
			Attribs:     make(map[string]string),
		}

		// Add any template variables to the check
		if val, ok := attribs[result.ID]; ok {
			// Decode each attribute document
			for _, doc := range val {
				err = json.Unmarshal(doc.Source, &result.Attribs)
				if err != nil {
					return nil, fmt.Errorf("Failed to decode attribute document from index %s for check %s: %s", doc.Index, doc.ID, err)
				}
			}
		}

		// Add the SavedValue attribute in case the check uses it
		result.Attribs["SavedValue"] = "{{.SavedValue}}"

		results = append(results, result)
	}

	return results, nil
}
