package checks

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
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
func RunChecks(defPass chan []schema.CheckDef, pubQueue chan<- beat.Event) {
	start := time.Now()

	// Recieve definitions from channel
	defs := <-defPass

	// Prepare event queue
	queue := make(chan schema.CheckResult, len(defs))

	// Iterate over each check
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	for _, def := range defs {
		checkName := def.Name
		check := unpackDef(def)

		// Start check goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()

			checkStart := time.Now()
			queue <- check.Run(ctx)
			logp.Info("[%s] Finished after %.2f seconds", checkName, time.Now().Since(checkStart).Seconds())
		}()
	}
	// Send definitions back through channel
	defPass <- defs

	// Wait for checks to finish
	wg.Wait()
	logp.Info("Checks started at %s have finished in %.2f seconds", start.Format("15:04:05.000"), time.Now().Since(start).Seconds())
	close(queue)
	for result := range queue {
		// Publish check results
		event := beat.Event{
			Timestamp: result.Timestamp,
			Fields: common.MapStr{
				"type":         "dynamicbeat",
				"id":           result.ID,
				"name":         result.Name,
				"group":        result.Group,
				"score_weight": result.ScoreWeight,
				"check_type":   result.CheckType,
				"passed":       result.Passed,
				"message":      result.Message,
				"details":      result.Details,
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
		fmt.Printf("\n\n[!] Add your definition to the switch case!\n\n")
	}
	err = def.Init(c.ID, c.Name, c.Group, c.ScoreWeight, renderedJSON)
	if err != nil {
		logp.Info("%s", err)
	}

	return def
}
