// Package cmdtask — task_audit.go provides centralized audit logging for macro, ssh, and installer subsystems.
package cmdtask

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RecordTaskAudit asynchronously or safely logs an operation into the split-db TaskHistory table.
func RecordTaskAudit(section, action, target, forwardPayload, status string) {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return
	}
	defer tasksDB.Close()

	taskID := fmt.Sprintf("%s-%d", section, time.Now().UnixNano())
	_ = tasksDB.InsertTaskHistory(taskID, section, action, target, forwardPayload, "", status)
}
