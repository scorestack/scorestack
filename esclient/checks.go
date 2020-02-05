package esclient

import (
	"strings"

	"github.com/elastic/go-elasticsearch"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/schema"
)

// UpdateCheckDefinitions : Re-read check definitions and attributes from Elasticsearch.
func UpdateCheckDefinitions(c *elasticsearch.Client, i string) (schema.CheckDefinitions, error) {
	var result schema.CheckDefinitions

	// Get list of checks
	checks, err := GetAllDocuments(c, i)
	if err != nil {
		return result, err
	}

	// Iterate over each check
	allAttribs := make(map[string]map[string]string, 0)
	for _, check := range checks {
		attribs := make(map[string]string)

		// Get any template variables for the check
		id := check["id"].String()
		for _, perm := range []string{"admin", "user"} {
			// Generate attribute index name
			idx := strings.Join([]string{"attrib_", perm, "_", id}, "")

			// Get attribute document
			docs, err := GetAllDocuments(c, idx)
			if err != nil {
				return result, err
			}
			attrib := docs[0]

			// Read attributes from document
			for k, v := range attrib {
				if _, pres := attribs[k]; !pres {
					attribs[k] = v.String()
				}
			}
		}

		// Add attributes to full map of attributes
		allAttribs[id] = attribs
	}

	result = schema.CheckDefinitions{
		Checks:     checks,
		Attributes: allAttribs,
	}
	return result, nil
}
