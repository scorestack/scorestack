package models

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ViewType: describes the user interface used to display the attribute's value
type ViewType string

const (
	DefaultView  ViewType = ""         // Equivalent to TextView
	BoolView     ViewType = "boolean"  // Toggle switch
	NumberView   ViewType = "number"   // Text box with numeric validation
	PasswordView ViewType = "password" // Censored text box with "view" button for sensitive values
	TextView     ViewType = "text"     // Standard text box (default)
)

// Permission: if teams can view the attribute and change its value
type Permission string

const (
	DefaultPerm Permission = ""       // Equivalent to NonePerm
	NonePerm    Permission = "none"   // Teams cannot read or modify the attribute or its value (default)
	ViewPerm    Permission = "view"   // Teams can view the attribute and its value
	UpdatePerm  Permission = "update" // Teams can view the attribute and change its value
)

type AttributeMetadata struct {
	CheckId     string     // Globally unique reference to the check this attribute applies to
	Group       string     // The team this attribute belongs to
	Key         string     // Key used in template rendering to reference the attribute
	Permissions Permission // Whether access should be restricted to the attribute, and how
}

func (a *AttributeMetadata) Id() string {
	return fmt.Sprintf("%s-%s", a.CheckId, a.Key)
}

// Attribute: a variable templated into a check definition that can be modified
// at runtime by competition organizers and participants. Attributes allow the
// behavior of checks to be modified mid-competition.
type Attribute struct {
	AttributeMetadata
	DisplayName string           // Human-readable name displayed in UIs
	Description string           // Explanation of what the attribute is for, displayed in UIs
	DisplayAs   ViewType         // How the attribute value should be rendered in UIs
	Values      []AttributeValue // All values the attribute has ever had
}

func (a *Attribute) Error(msg string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s for attribute check_id='%s' key='%s': %s", msg, a.CheckId, a.Key, err)
	}

	return fmt.Sprintf("%s for attribute check_id='%s' key='%s'", msg, a.CheckId, a.Key)
}

// Value: get the current value string for the attribute. Returns an empty
// string if no values have been defined for the attribute yet.
func (a *Attribute) Value() string {
	var value string
	var newestTime time.Time

	for _, val := range a.Values {
		if val.Created.After(newestTime) {
			newestTime = val.Created
			value = val.Value
		}
	}

	// Since attributes should always be created with an initial value,
	// something went wrong if the attribute doesn't have any values.
	if newestTime.IsZero() {
		zap.S().Errorf(a.Error("no value found", nil))
	}

	return value
}

// ---[esclient.Indexable]-----------------------------------------------------

func (a *Attribute) GetIndex() string {
	switch a.Permissions {
	case UpdatePerm:
		return fmt.Sprintf("attributes-%s", a.Group)
	case ViewPerm:
		return fmt.Sprintf("attributes-view-%s", a.Group)
	default:
		return "attributes-admin"
	}
}

func (a *Attribute) GetDocumentId() string {
	return a.Id()
}

func (a *Attribute) GetBody() map[string]interface{} {
	return map[string]interface{}{
		"check_id":     a.CheckId,
		"key":          a.Key,
		"display_name": a.DisplayName,
		"description":  a.Description,
		"display_as":   string(a.DisplayAs),
	}
}
