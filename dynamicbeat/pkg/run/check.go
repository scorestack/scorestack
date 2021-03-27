package run

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"text/template"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes"
	"go.uber.org/zap"
)

func Check(ctx context.Context, def check.Config) check.Result {
	// Create a check from the definition
	chk, err := unpackDef(def)
	if err != nil {
		return check.Result{
			Timestamp: time.Now(),
			Metadata:  def.Metadata,
			Passed:    false,
			Message:   fmt.Sprintf("encountered an error when unpacking check definition: %s", err),
			Details:   nil,
		}
	}

	// Set up the channel to recieve the CheckResult from the Check
	result := make(chan check.Result, 1)

	// Run the check
	go func() {
		result <- chk.Run(ctx)
	}()

	// Wait for either the timeout or for the check to finish
	for {
		select {
		case <-ctx.Done():
			// We already initialized the event with the correct values for a
			// context timeout, so just return that.
			return check.Result{
				Timestamp: time.Now(),
				Metadata:  def.Metadata,
				Passed:    false,
				Message:   "check timed out",
				Details:   nil,
			}
		case r := <-result:
			close(result)
			return r
		}
	}
}

func unpackDef(config check.Config) (check.Check, error) {
	// Render any template strings in the definition
	var renderedJSON []byte
	templ := template.New("definition")
	templ, err := templ.Parse(string(config.Definition))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse template for check: %s", err.Error())
	}

	var buf bytes.Buffer
	err = templ.Execute(&buf, config.Attributes.Merged())
	if err != nil {
		return nil, fmt.Errorf("Failed to execute template for check: %s", err.Error())
	}

	renderedJSON = buf.Bytes()

	// Create a Definition from the rendered JSON string
	def := checktypes.GetCheckType(config)
	err = initCheck(config, renderedJSON, def)
	if err != nil {
		zap.S().Infof("%s", err)
	}

	return def, nil
}

func initCheck(config check.Config, def []byte, chk check.Check) error {
	// Unpack definition JSON
	err := json.Unmarshal(def, &chk)
	if err != nil {
		return err
	}

	// Set generic values
	chk.SetConfig(config)

	// Process the field options
	return processFields(chk, chk.GetConfig().ID, chk.GetConfig().Type)
}

func processFields(s interface{}, id string, typ string) error {
	// Convert the parameter to reflect.Type and reflect.Value variables
	fields := reflect.TypeOf(s)
	if fields.Kind() == reflect.Ptr {
		fields = fields.Elem()
	}
	values := reflect.ValueOf(s)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	// Process each field in the struct
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)
		optiontype := field.Tag.Get("optiontype")

		switch optiontype {
		case "required":
			// Make sure the value is nonzero
			if value.IsZero() {
				return check.ValidationError{
					ID:    id,
					Type:  typ,
					Field: field.Name,
				}
			}
		case "optional":
			dflt := field.Tag.Get("optiondefault")

			// If the optiondefault is not set, then don't do anything with
			// this field. This typically means that the default for the field
			// is the zero value for the type, in which case we don't have to
			// do anything else.
			if dflt == "" {
				continue
			}

			// If the value is still zero, set the default value
			if value.IsZero() {
				switch value.Kind() {
				case reflect.Bool:
					v, _ := strconv.ParseBool(dflt)
					value.SetBool(v)
				case reflect.Float32, reflect.Float64:
					v, _ := strconv.ParseFloat(dflt, 64)
					value.SetFloat(v)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, _ := strconv.ParseInt(dflt, 0, 64)
					value.SetInt(v)
				case reflect.String:
					value.SetString(dflt)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v, _ := strconv.ParseUint(dflt, 0, 64)
					value.SetUint(v)
				}
			}
		case "list":
			// Recurse on each item in the list
			for j := 0; j < value.Len(); j++ {
				err := processFields(value.Index(j).Interface(), id, typ)
				if err != nil {
					return err
				}
			}
		default:
			// If the optiontype is invalid, or no optiontype is set, then don't do anything with this field
		}
	}

	return nil
}
