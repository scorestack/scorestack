package run

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checktypes"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/event"
	"go.uber.org/zap"
)

// RunChecks : Run a course of checks based on the currently-loaded configuration.
func RunChecks(defPass chan []check.Config, pubQueue chan<- event.Event) {
	start := time.Now()

	// Recieve definitions from channel
	defs := <-defPass
	zap.S().Infof("Recieved defs")

	// Make an event queue separate from the publisher queue so we can track
	// which checks are still running
	eventQueue := make(chan event.Event, len(defs))

	// Iterate over each check
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	names := make(map[string]bool)
	var wg sync.WaitGroup
	for _, def := range defs {
		// checkName := def.Name
		names[def.Meta.ID] = false
		chk, err := unpackDef(def)
		if err != nil {
			// Something was wrong with templating the check. Return a failed event with the error.
			errorDetail := make(map[string]string)
			errorDetail["error_message"] = err.Error()
			eventQueue <- event.Event{
				Timestamp:   time.Now(),
				Id:          chk.GetConfig().Meta.ID,
				Name:        chk.GetConfig().Meta.Name,
				CheckType:   chk.GetConfig().Meta.Type,
				Group:       chk.GetConfig().Meta.Group,
				ScoreWeight: chk.GetConfig().Meta.ScoreWeight,
				Passed:      false,
				Message:     "Encountered an error when unpacking check definition.",
				Details:     errorDetail,
			}
		}

		// Start check goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()

			checkStart := time.Now()
			checkName := chk.GetConfig().Meta.Name
			eventQueue <- runCheck(ctx, chk)
			zap.S().Infof("[%s] Finished after %.2f seconds", checkName, time.Since(checkStart).Seconds())
		}()
	}
	// Send definitions back through channel
	defPass <- defs

	// Wait for checks to finish
	defer wg.Wait()
	zap.S().Infof("Checks started at %s have finished in %.2f seconds", start.Format("15:04:05.000"), time.Since(start).Seconds())
	go func() {
		for {
			if names == nil {
				break
			} else if len(names) == 0 {
				break
			} else {
				time.Sleep(30 * time.Second)
				zap.S().Infof("Checks still running after %.2f seconds: %+v", time.Since(start).Seconds(), names)
			}
		}
		zap.S().Infof("All checks started %.2f seconds ago have finished", time.Since(start).Seconds())
		close(eventQueue)
	}()
	for evt := range eventQueue {
		// Record that the check has finished
		delete(names, evt.Id)

		// Publish the event to the publisher queue
		pubQueue <- evt
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
	err = templ.Execute(&buf, config.Attribs)
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

func runCheck(ctx context.Context, chk check.Check) event.Event {
	// Initialize the event to be published
	evt := event.Event{
		Timestamp:   time.Now(),
		Id:          chk.GetConfig().Meta.ID,
		Name:        chk.GetConfig().Meta.Name,
		CheckType:   chk.GetConfig().Meta.Type,
		Group:       chk.GetConfig().Meta.Group,
		ScoreWeight: chk.GetConfig().Meta.ScoreWeight,
		Passed:      false,
		Message:     "Check timed out",
		Details:     nil,
	}

	// Set up the channel to recieve the CheckResult from the Check
	recieveResult := make(chan check.Result, 1)

	// Run the check
	go func() {
		recieveResult <- chk.Run(ctx)
	}()

	// Wait for either the timeout or for the check to finish
	for {
		select {
		case <-ctx.Done():
			// We already initialized the event with the correct values for a
			// context timeout, so just return that.
			return evt
		case result := <-recieveResult:
			close(recieveResult)
			// Set the passed, message, and details fields with the CheckResult
			evt.Passed = result.Passed
			evt.Message = result.Message
			evt.Details = result.Details
			return evt
		}
	}
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
	return processFields(chk, chk.GetConfig().Meta.ID, chk.GetConfig().Meta.Type)
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
