package noop

import (
	"strings"
	"time"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/common"
)

// Run : Execute the check
func Run(chk common.Check) {
	defer chk.WaitGroup.Done()

	// Render the message string
	message := make([]string, 0)
	for k, v := range chk.Definition {
		message = append(message, k)
		message = append(message, ":")
		message = append(message, v)
		message = append(message, ",")
	}

	result := common.CheckResult{
		Timestamp: time.Now(),
		ID:        chk.ID,
		Name:      chk.Name,
		CheckType: "noop",
		Passed:    true,
		Message:   strings.Join(message, " "),
		Details:   chk.Definition,
	}

	chk.Output <- result
}
