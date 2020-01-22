package checks

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/elastic/beats/libbeat/beat"
	beatcommon "github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/common"
	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/noop"
)

// RunChecks : Run a course of checks based on the currently-loaded configuration.
func RunChecks(client beat.Client, defs common.CheckDefinitions) {
	queue := make(chan common.CheckResult, len(defs.Checks))
	var wg sync.WaitGroup

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
			WaitGroup:  &wg,
			Output:     queue,
		}

		// Start check goroutine
		wg.Add(1)
		switch chk["type"].String() {
		case "noop":
			go noop.Run(chkInfo)
		}

		// Wait for checks to finish
		wg.Wait()
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
			client.Publish(event)
			logp.Info("Event sent")
		}
	}
}
