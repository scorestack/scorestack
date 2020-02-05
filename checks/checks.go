package checks

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/http"
	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/noop"
	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/schema"
)

// RunChecks : Run a course of checks based on the currently-loaded configuration.
func RunChecks(defPass chan []schema.CheckDef, wg *sync.WaitGroup, pubQueue chan<- beat.Event) {
	defer wg.Done()

	// Recieve definitions from channel
	defs := <-defPass

	// Prepare event queue
	queue := make(chan schema.CheckResult, len(defs))
	var events sync.WaitGroup

	// Iterate over each check
	for _, def := range defs {
		check := unpackDef(def)

		// Construct Check struct
		chkInfo := schema.Check{
			ID:        chk["id"].String(),
			Name:      chk["name"].String(),
			WaitGroup: &events,
			Output:    queue,
		}

		// Add definitions to correct attribute in Check struct
		if chk["definition"].IsArray() {
			chkInfo.DefinitionList = defs
		} else {
			chkInfo.Definition = defs[0]
		}

		// Start check goroutine
		events.Add(1)
		switch chk["type"].String() {
		case "noop":
			go noop.Run(chkInfo)
		case "http":
			go http.Run(chkInfo)
		default:
			// We didn't start a goroutine, so the WaitGroup counter needs to be decremented.
			// If this wasn't here, events.Wait() would hang forever if there was a check with an unknown type.
			// This also allows us to have only one events.Add(1) at the beginning of the switch/case.
			// Otherwise, we would have to add a events.Add(1) to each case.
			events.Done()
		}

	}
	// Send definitions back through channel
	defPass <- defs

	// Wait for checks to finish
	events.Wait()
	close(queue)
	for result := range queue {
		// Publish check results
		event := beat.Event{
			Timestamp: result.Timestamp,
			Fields: common.MapStr{
				"type":       "dynamicbeat",
				"id":         result.ID,
				"name":       result.Name,
				"check_type": result.CheckType,
				"passed":     result.Passed,
				"message":    result.Message,
				"details":    result.Details,
			},
		}
		pubQueue <- event
	}
}

func unpackDef(c schema.CheckDef) schema.Check {
	// Render any template strings in the definition
	var renderedJSON []byte
	templ := template.Must(template.New("definition").Parse(string(c.Definition)))
	var buf bytes.Buffer
	err := templ.Execute(&buf, c.Attribs)
	if err != nil {
		// If there was an error parsing the template, use the original string
		renderedJSON = c.Definition
	} else {
		renderedJSON = buf.Bytes()
	}

	// Create a Definition from the rendered JSON string
	var def schema.Check
	switch c.Type {
	case "noop":
		def = noop.Definition{}
		def.Init(c.ID, c.Name, renderedJSON)
	case "http":
		def = http.Definition{}
		def.Init(c.ID, c.Name, renderedJSON)
	default:
	}

	return def
}
