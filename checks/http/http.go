package http

import (
	"time"

	"gitlab.ritsec.cloud/newman/dynamicbeat/checks/common"
)

// Run : Execute the check
func Run(chk common.Check) {
	defer chk.WaitGroup.Done()

	result := common.CheckResult{
		Timestamp: time.Now(),
		ID:        chk.ID,
		Name:      chk.Name,
		CheckType: "http",
		Passed:    true,
		Message:   "TODO",
		Details:   chk.Definition,
	}

	chk.Output <- result
}
