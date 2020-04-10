package checks

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/s-newman/scorestack/dynamicbeat/checks/dns"
	"github.com/s-newman/scorestack/dynamicbeat/checks/ftp"
	"github.com/s-newman/scorestack/dynamicbeat/checks/http"
	"github.com/s-newman/scorestack/dynamicbeat/checks/icmp"
	"github.com/s-newman/scorestack/dynamicbeat/checks/imap"
	"github.com/s-newman/scorestack/dynamicbeat/checks/ldap"
	"github.com/s-newman/scorestack/dynamicbeat/checks/noop"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
	"github.com/s-newman/scorestack/dynamicbeat/checks/smtp"
	"github.com/s-newman/scorestack/dynamicbeat/checks/ssh"
	"github.com/s-newman/scorestack/dynamicbeat/checks/vnc"
	"github.com/s-newman/scorestack/dynamicbeat/checks/winrm"
	"github.com/s-newman/scorestack/dynamicbeat/checks/xmpp"
)

// RunChecks : Run a course of checks based on the currently-loaded configuration.
func RunChecks(defPass chan []schema.CheckConfig, pubQueue chan<- beat.Event) {
	start := time.Now()

	// Recieve definitions from channel
	defs := <-defPass
	logp.Info("Recieved defs")

	// Make an event queue separate from the publisher queue so we can track
	// which checks are still running
	eventQueue := make(chan beat.Event, len(defs))

	// Iterate over each check
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	names := make(map[string]bool)
	var wg sync.WaitGroup
	for _, def := range defs {
		// checkName := def.Name
		names[def.ID] = false
		check := unpackDef(def)

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
		logp.Info("All checks started %.2f seconds ago have finished", time.Since(start).Seconds())
		close(eventQueue)
	}()
	for event := range eventQueue {
		// Record that the check has finished
		id, _ := event.Fields.GetValue("id")
		delete(names, id.(string))

		// Publish the event to the publisher queue
		pubQueue <- event
	}
}

func unpackDef(config schema.CheckConfig) schema.Check {
	// Render any template strings in the definition
	var renderedJSON []byte
	templ := template.Must(template.New("definition").Parse(string(config.Definition)))
	var buf bytes.Buffer
	err := templ.Execute(&buf, config.Attribs)
	if err != nil {
		// If there was an error parsing the template, use the original string
		renderedJSON = config.Definition
	} else {
		renderedJSON = buf.Bytes()
	}

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
	default:
		logp.Warn("Invalid check type found. Offending check : %s:%s", config.Name, config.Type)
		def = &noop.Definition{}
	}
	err = initCheck(config, renderedJSON, def)
	if err != nil {
		logp.Info("%s", err)
	}

	return def
}

func runCheck(ctx context.Context, check schema.Check) beat.Event {
	// Initialize the event to be published
	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":         "dynamicbeat",
			"id":           check.GetConfig().ID,
			"name":         check.GetConfig().Name,
			"check_type":   check.GetConfig().Type,
			"group":        check.GetConfig().Group,
			"score_weight": check.GetConfig().ScoreWeight,
			"passed":       false,
			"message":      "Check timed out",
			"details":      nil,
		},
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
			return event
		case result := <-recieveResult:
			close(recieveResult)
			// Set the passed, message, and details fields with the CheckResult
			event.Fields.Put("passed", result.Passed)
			event.Fields.Put("message", result.Message)
			event.Fields.Put("details", result.Details)
			return event
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
	fields := reflect.TypeOf(s).Elem()
	values := reflect.ValueOf(s).Elem()

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
