package checks

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/elastic/beats/libbeat/beat"
	beatcommon "github.com/elastic/beats/libbeat/common"
	"github.com/tidwall/gjson"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/common"
	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/http"
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
		defs := unpackDefs(chk, defs.Attributes)

		// Construct Check struct
		chkInfo := common.Check{
			ID:        chk["id"].String(),
			Name:      chk["name"].String(),
			WaitGroup: &events,
			Output:    queue,
		}

		// Add definitions to correct attribute in Check struct
		if len(defs) > 1 {
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

func unpackDefs(check map[string]gjson.Result, attribs map[string]map[string]string) []map[string]string {
	// The definition can be an array, so we assume it is an array.
	// If the definition is just a map, create an array of length 1 with it.
	var packedDefs []gjson.Result
	if check["definition"].IsArray() {
		packedDefs = check["definition"].Array()
	} else {
		packedDefs = []gjson.Result{check["definition"]}
	}

	// Unpack each definition
	unpackedDefs := make([]map[string]string, 0)
	for _, packedDef := range packedDefs {
		// Template out the contents of the definition
		def := make(map[string]string)
		packedMap := packedDef.Map()
		for k, v := range packedMap {
			// Render template string in value, if any
			templ := template.Must(template.New(k).Parse(v.String()))
			var buf bytes.Buffer
			templ.Execute(&buf, attribs[check["id"].String()])
			def[k] = buf.String()
		}
		unpackedDefs = append(unpackedDefs, def)
	}

	return unpackedDefs
}
