package models

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// TODO
type CheckConfig struct {
	CheckMetadata
	Definition []byte // The untemplated check definition as a JSON string
	Attributes []Attribute
}

func (c *CheckConfig) Error(msg string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s for check config check_id='%s' check_kind='%s': %s", msg, c.CheckId, c.Kind, err)
	}

	return fmt.Sprintf("%s for check config check_id='%s' check_kind='%s'", msg, c.CheckId, c.Kind)
}

func (c *CheckConfig) MergedAttributes() map[string]string {
	attributes := make(map[string]string)

	for _, attribute := range c.Attributes {
		attributes[attribute.Key] = attribute.Value()
	}

	return attributes
}

// ---[esclient.Indexable]-----------------------------------------------------

func (c *CheckConfig) GetIndex() string {
	return "checks"
}

func (c *CheckConfig) GetDocumentId() string {
	// This is redundant but I prefer to make it explicit that CheckConfig
	// implements the Indexable interface
	return c.CheckMetadata.GetDocumentId()
}

func (c *CheckConfig) GetBody() map[string]interface{} {
	definition := make(map[string]interface{})
	err := json.Unmarshal(c.Definition, &definition)
	if err != nil {
		// The Indexable interface doesn't let us return an error from this
		// function, so we just have to log and continue. The definition will
		// be an empty object.
		zap.S().Error(c.Error("failed to marshal definition from JSON string", err))
	}

	document := c.CheckMetadata.GetBody()
	document["definition"] = definition

	return document
}
