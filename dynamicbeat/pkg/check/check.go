package check

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Check interface {
	GetConfig() Config
	SetConfig(c Config)
	Run(ctx context.Context) Result
}

type Metadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Group       string `json:"group"`
	ScoreWeight int64  `json:"score_weight"`
}

type Config struct {
	Metadata
	Definition []byte
	Attributes `json:"attributes"`
}

type AttributeType string

const (
	ATTRIBUTE_BOOLEAN  AttributeType = "boolean"
	ATTRIBUTE_NUMBER   AttributeType = "number"
	ATTRIBUTE_PASSWORD AttributeType = "password"
	ATTRIBUTE_TEXT     AttributeType = "text"
)

type Attribute struct {
	CheckID      string        `json:"check_id"`
	Name         string        `json:"name"`
	InitialValue string        `json:"initial_value"`
	DisplayName  string        `json:"display_name"`
	Description  string        `json:"description"`
	Type         AttributeType `json:"type"`
}

type AttributeValue struct {
	CheckID string    `json:"check_id"`
	Name    string    `json:"name"`
	Value   string    `json:"value"`
	Created time.Time `json:"created"`
}

type Attributes struct {
	Admin map[string]Attribute `json:"admin"`
	User  map[string]Attribute `json:"user"`
}

type IndexDoc struct {
	Index string
	Doc   io.Reader
}

func (a *Attributes) Merged() map[string]Attribute {
	m := make(map[string]Attribute)

	for k, v := range a.Admin {
		m[k] = v
	}
	for k, v := range a.User {
		m[k] = v
	}

	return m
}

// A ValidationError represents an issue with a check definition.
type ValidationError struct {
	ID    string // the ID of the check with an invalid definition
	Type  string // the type of the check with an invalid definition
	Field string // the field in the check definition that was invalid
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Error: check (Type: `%s`, ID: `%s`) is missing value for required field `%s`", v.Type, v.ID, v.Field)
}

func (c *Config) Documents() (io.Reader, io.Reader, []IndexDoc, error) {
	def := make(map[string]interface{})
	err := json.Unmarshal(c.Definition, &def)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal definition for '%s': %s", c.ID, err)
	}

	// The check definition document doesn't include the attributes
	chk := struct {
		Metadata
		Definition map[string]interface{} `json:"definition"`
	}{c.Metadata, def}
	checkDoc, err := json.Marshal(chk)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal definition for '%s': %s", c.ID, err)
	}

	// The generic check definition only includes the metadata
	generic := struct {
		Metadata
	}{c.Metadata}
	genericDoc, err := json.Marshal(generic)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal generic definition for '%s': %s", c.ID, err)
	}

	attributes := make([]IndexDoc, 0)

	err = attributeDocs(c.Attributes.Admin, &attributes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal admin attributes for '%s': %s", c.ID, err)
	}

	err = attributeDocs(c.Attributes.User, &attributes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal user attributes for '%s': %s", c.ID, err)
	}

	return bytes.NewReader(checkDoc), bytes.NewReader(genericDoc), attributes, nil
}

func attributeDocs(attributes map[string]Attribute) ([]error {
	if attributes == nil {
		return nil
	}

	for name, attrib := range attributes {
		
	}

	doc, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(doc), nil
}
