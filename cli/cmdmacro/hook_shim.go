package cmdmacro

import (
	"time"
)

// FinishCommandAuditFn is wired by cmd/di_hooks.go to the canonical implementation.
var FinishCommandAuditFn func(shouldAudit bool, id int64, start time.Time, exitCode int, summary string, repoCount int)

func finishCommandAudit(shouldAudit bool, id int64, start time.Time, exitCode int, summary string, repoCount int) {
	if FinishCommandAuditFn != nil {
		FinishCommandAuditFn(shouldAudit, id, start, exitCode, summary, repoCount)
	}
}
