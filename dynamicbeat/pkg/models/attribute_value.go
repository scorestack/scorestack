package models

import (
	"fmt"
	"time"
)

// AttributeValue: the raw string that will be templated into the check
// definition where the attribute is referenced. An attribute can have multiple
// values, but only the most recent value (as determined by the attribute's
// created time) will be templated into the check definition.
type AttributeValue struct {
	AttributeMetadata
	Value   string    // Value inserted into check definition when referenced
	Created time.Time // System time of when the attribute was first created
}

// ---[esclient.Indexable]-----------------------------------------------------

func (a *AttributeValue) GetIndex() string {
	switch a.Permissions {
	case UpdatePerm:
		return fmt.Sprintf("attribute-values-%s", a.Group)
	case ViewPerm:
		return fmt.Sprintf("attribute-values-view-%s", a.Group)
	default:
		return "attribute-values-admin"
	}
}

func (a *AttributeValue) GetDocumentId() string {
	return ""
}

func (a *AttributeValue) GetBody() map[string]interface{} {
	return map[string]interface{}{
		"check_id": a.CheckId,
		"key":      a.Key,
		"value":    a.Value,
		"created":  a.Created.Format(time.RFC3339),
	}
}
