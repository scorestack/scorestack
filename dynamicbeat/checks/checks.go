package checks

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"

	"github.com/s-newman/scorestack/dynamicbeat/checks/http"
	"github.com/s-newman/scorestack/dynamicbeat/checks/icmp"
	"github.com/s-newman/scorestack/dynamicbeat/checks/noop"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
	"github.com/s-newman/scorestack/dynamicbeat/checks/ssh"
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

		// Start check goroutine
		events.Add(1)
		check.Run(&events, queue)
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
		def = &noop.Definition{}
		def.Init(c.ID, c.Name, renderedJSON)
	case "http":
		def = &http.Definition{}
		def.Init(c.ID, c.Name, renderedJSON)
	case "icmp":
		def = &icmp.Definition{}
		def.Init(c.ID, c.Name, renderedJSON)
	case "ssh":
		def = &ssh.Definition{}
		def.Init(c.ID, c.Name, renderedJSON)
	default:
		fmt.Printf("Add your definition to the switch case!\n")
	}

	return def
}
