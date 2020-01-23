package checks

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/elastic/beats/libbeat/beat"
	beatcommon "github.com/elastic/beats/libbeat/common"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/common"
	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/noop"
)

// RunChecks : Run a course of checks based on the currently-loaded configuration.
func RunChecks(defPass chan common.CheckDefinitions, wg *sync.WaitGroup, pubQueue chan<- beat.Event) {
	defer wg.Done()

	// Recieve definitions from channel
	defs := <-defPass

	// Prepare event queue
	queue := make(chan common.CheckResult, len(defs.Checks))
	var events sync.WaitGroup

	// Iterate over each check
	for _, chk := range defs.Checks {
		// Unpack definition
		packedDef := chk["definition"].Map()
		def := make(map[string]string)
		for k, v := range packedDef {
			// Render template string in value, if any
			templ := template.Must(template.New(k).Parse(v.String()))
			var buf bytes.Buffer
			if err := templ.Execute(&buf, defs.Attributes[chk["id"].String()]); err != nil {
				// TODO: pass error back through channel
			}
			def[k] = buf.String()
		}

		// Construct Check struct
		chkInfo := common.Check{
			ID:         chk["id"].String(),
			Name:       chk["name"].String(),
			Definition: def,
			WaitGroup:  &events,
			Output:     queue,
		}

		// Start check goroutine
		events.Add(1)
		switch chk["type"].String() {
		case "noop":
			go noop.Run(chkInfo)
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
			Fields: beatcommon.MapStr{
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
