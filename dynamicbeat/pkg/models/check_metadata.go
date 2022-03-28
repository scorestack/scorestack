package models

// CheckMetadata: configuration fields common to all check types
type CheckMetadata struct {
	CheckId     string // Globally unique identifier for the check
	Group       string // The team this check belongs to
	DisplayName string // Human-readable name displayed in UIs
	Description string // Explanation of what the check does, displayed in UIs
	Kind        string // The kind of check that will be run
	Points      uint64 // How many points a team is awarded for a passing check
}

// ---[esclient.Indexable]-----------------------------------------------------

func (c *CheckMetadata) GetIndex() string {
	return "check-views"
}

func (c *CheckMetadata) GetDocumentId() string {
	return c.CheckId
}

func (c *CheckMetadata) GetBody() map[string]interface{} {
	return map[string]interface{}{
		"check_id":     c.CheckId,
		"group":        c.Group,
		"display_name": c.DisplayName,
		"description":  c.Description,
		"kind":         c.Kind,
		"points":       c.Points,
	}
}
