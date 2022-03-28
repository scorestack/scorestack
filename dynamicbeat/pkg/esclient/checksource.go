package esclient

import (
	"encoding/json"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/models"
	"go.uber.org/zap"
)

func (c *Client) LoadAll() ([]models.CheckConfig, error) {
	// Track how long it takes to update check definitions
	start := time.Now()

	checks, err := c.GetChecks()
	if err != nil {
		return nil, err
	}
	attributes, err := c.GetAttributes()
	if err != nil {
		return nil, err
	}

	for _, attribute := range attributes {
		if check, exists := checks[attribute.CheckId]; exists {
			check.Attributes = append(check.Attributes, attribute)
		} else {
			zap.S().Errorf("ignoring attribute key=%s defined for non-existant check check_id=%s", attribute.Key, attribute.CheckId)
			continue
		}
	}

	out := make([]models.CheckConfig, 0)
	for _, check := range checks {
		out = append(out, check)
	}

	elaspsed := time.Since(start)
	zap.S().Infof("loaded %d checks and %d attributes from Elasticsearch in %s", len(checks), len(attributes), elaspsed)
	return out, nil
}

// GetChecks: load all check documents from Elasticsearch. Attributes for each
// check must be populated separately. Checks are returned as a map of check ID
// strings to CheckConfig structs.
func (c *Client) GetChecks() (map[string]models.CheckConfig, error) {
	docs, err := c.GetAllDocumentsFrom("checks-*")
	if err != nil {
		return nil, err
	}

	checks := make(map[string]models.CheckConfig, 0)
	for _, doc := range docs {
		check := models.CheckConfig{
			CheckMetadata: models.CheckMetadata{
				CheckId:     doc.TryLoad("check_id"),
				Group:       doc.TryLoad("group"),
				DisplayName: doc.TryLoad("display_name"),
				Description: doc.TryLoad("description"),
				Kind:        doc.TryLoad("kind"),
			},
			Attributes: make([]models.Attribute, 0),
		}

		// Try to load the points string
		if points, ok := doc.Source["points"].(uint64); ok {
			check.Points = points
		} else {
			zap.S().Errorf("ignoring check check_id=%s due to invalid points field", check.CheckId)
			continue
		}

		// Try to convert the definition back to JSON
		buf, err := json.Marshal(doc.Source["definition"])
		if err != nil {
			zap.S().Errorf("ignoring check check_id=%s due to error marshalling definition: %s", check.CheckId, err)
		}
		check.Definition = buf

		checks[check.CheckId] = check
	}

	return checks, nil
}

// GetAttributes: load all attributes and their values from Elasticsearch.
func (c *Client) GetAttributes() ([]models.Attribute, error) {
	values, err := c.GetAttributeValues()
	if err != nil {
		return nil, err
	}

	docs, err := c.GetAllDocumentsFrom("attributes-*")
	if err != nil {
		return nil, err
	}

	attributes := make([]models.Attribute, 0)
	for _, doc := range docs {
		attribute := models.Attribute{
			AttributeMetadata: models.AttributeMetadata{
				CheckId: doc.TryLoad("check_id"),
				Key:     doc.TryLoad("key"),
				// We don't need the Group or Permissions fields to run the
				// check so they're left empty
			},
			DisplayName: doc.Source["display_name"].(string),
			Description: doc.Source["description"].(string),
			DisplayAs:   doc.Source["display_as"].(models.ViewType),
			Values:      make([]models.AttributeValue, 0),
		}

		if valueList, exists := values[attribute.Id()]; exists {
			for _, value := range valueList {
				attribute.Values = append(attribute.Values, value)
			}
		}

		attributes = append(attributes, attribute)
	}

	return attributes, nil
}

func (c *Client) GetAttributeValues() (map[string][]models.AttributeValue, error) {
	docs, err := c.GetAllDocumentsFrom("attribute-values-*")
	if err != nil {
		return nil, err
	}

	values := make(map[string][]models.AttributeValue)
	for _, doc := range docs {
		value := models.AttributeValue{
			AttributeMetadata: models.AttributeMetadata{
				CheckId: doc.TryLoad("check_id"),
				Key:     doc.TryLoad("key"),
				// We don't need the Group or Permissions fields to run the
				// check so they're left empty
			},
			Value: doc.TryLoad("value"),
		}

		created, err := time.Parse(time.RFC3339, doc.TryLoad("created"))
		if err != nil {
			zap.S().Errorf("skipping attribute value in document_id=%s index=%s due to issue parsing creation timestamp: %s", doc.ID, doc.Index, err)
			continue
		}
		value.Created = created

		if _, exists := values[value.Id()]; exists {
			values[value.Id()] = append(values[value.Id()], value)
		} else {
			values[value.Id()] = []models.AttributeValue{value}
		}
	}

	return values, nil
}
