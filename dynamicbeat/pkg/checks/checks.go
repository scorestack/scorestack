package checks

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

	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/dns"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/ftp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/http"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/icmp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/imap"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/ldap"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/mssql"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/mysql"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/noop"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/postgresql"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/schema"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/smb"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/smtp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/ssh"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/vnc"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/winrm"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/checks/xmpp"
	"github.com/scorestack/scorestack/dynamicbeat/pkg/event"
)

// RunChecks : Run a course of checks based on the currently-loaded configuration.
func RunChecks(defPass chan []schema.CheckConfig, pubQueue chan<- event.Event) {
	// start := time.Now()

	// Recieve definitions from channel
	defs := <-defPass
	// logp.Info("Recieved defs")

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
		names[def.ID] = false
		check, err := unpackDef(def)
		if err != nil {
			// Something was wrong with templating the check. Return a failed event with the error.
			errorDetail := make(map[string]string)
			errorDetail["error_message"] = err.Error()
			eventQueue <- event.Event{
				Timestamp:   time.Now(),
				Id:          check.GetConfig().ID,
				Name:        check.GetConfig().Name,
				CheckType:   check.GetConfig().Type,
				Group:       check.GetConfig().Group,
				ScoreWeight: check.GetConfig().ScoreWeight,
				Passed:      false,
				Message:     "Encountered an error when unpacking check definition.",
				Details:     errorDetail,
			}
		}

		// Start check goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()

			// checkStart := time.Now()
			eventQueue <- runCheck(ctx, check)
			// logp.Info("[%s] Finished after %.2f seconds", checkName, time.Since(checkStart).Seconds())
		}()
	}
	// Send definitions back through channel
	defPass <- defs

	// Wait for checks to finish
	defer wg.Wait()
	// logp.Info("Checks started at %s have finished in %.2f seconds", start.Format("15:04:05.000"), time.Since(start).Seconds())
	go func() {
		for {
			if names == nil {
				break
			} else if len(names) == 0 {
				break
			} else {
				time.Sleep(30 * time.Second)
				// logp.Info("Checks still running after %.2f seconds: %+v", time.Since(start).Seconds(), names)
			}
		}
		// logp.Info("All checks started %.2f seconds ago have finished", time.Since(start).Seconds())
		close(eventQueue)
	}()
	for evt := range eventQueue {
		// Record that the check has finished
		delete(names, evt.Id)

		// Publish the event to the publisher queue
		pubQueue <- evt
	}
}

func unpackDef(config schema.CheckConfig) (schema.Check, error) {
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
	var def schema.Check
	switch config.Type {
	case "noop":
		def = &noop.Definition{}
	case "http":
		def = &http.Definition{}
	case "icmp":
		def = &icmp.Definition{}
	case "ssh":
		def = &ssh.Definition{}
	case "dns":
		def = &dns.Definition{}
	case "ftp":
		def = &ftp.Definition{}
	case "ldap":
		def = &ldap.Definition{}
	case "vnc":
		def = &vnc.Definition{}
	case "imap":
		def = &imap.Definition{}
	case "smtp":
		def = &smtp.Definition{}
	case "winrm":
		def = &winrm.Definition{}
	case "xmpp":
		def = &xmpp.Definition{}
	case "mysql":
		def = &mysql.Definition{}
	case "smb":
		def = &smb.Definition{}
	case "postgresql":
		def = &postgresql.Definition{}
	case "mssql":
		def = &mssql.Definition{}
	default:
		// logp.Warn("Invalid check type found. Offending check : %s:%s", config.Name, config.Type)
		def = &noop.Definition{}
	}
	err = initCheck(config, renderedJSON, def)
	if err != nil {
		// logp.Info("%s", err)
	}

	return def, nil
}

func runCheck(ctx context.Context, check schema.Check) event.Event {
	// Initialize the event to be published
	evt := event.Event{
		Timestamp:   time.Now(),
		Id:          check.GetConfig().ID,
		Name:        check.GetConfig().Name,
		CheckType:   check.GetConfig().Type,
		Group:       check.GetConfig().Group,
		ScoreWeight: check.GetConfig().ScoreWeight,
		Passed:      false,
		Message:     "Check timed out",
		Details:     nil,
	}

	// Set up the channel to recieve the CheckResult from the Check
	recieveResult := make(chan schema.CheckResult, 1)

	// Run the check
	go func() {
		recieveResult <- check.Run(ctx)
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

func initCheck(config schema.CheckConfig, def []byte, check schema.Check) error {
	// Unpack definition JSON
	err := json.Unmarshal(def, &check)
	if err != nil {
		return err
	}

	// Set generic values
	check.SetConfig(config)

	// Process the field options
	return processFields(check, check.GetConfig().ID, check.GetConfig().Type)
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
				return schema.ValidationError{
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
