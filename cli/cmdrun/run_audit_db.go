package cmdrun

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// EnqueueRunTaskAudit inserts a running execution task into the TaskQueue.
func EnqueueRunTaskAudit(target *RunTarget, args []string) (string, error) {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return "", apperror.WrapSimple(err, "cmdrun.enqueue_audit: open tasks db")
	}
	defer tasksDB.Close()

	taskId := fmt.Sprintf("run-%d", time.Now().UnixNano())
	payload := strings.Join(args, " ")
	query := `INSERT INTO TaskQueue (QueueId, Section, Action, Target, ForwardPayload, Status) VALUES (?, 'run', 'file-run', ?, ?, 'running')`

	res := store.ExecWrapper(tasksDB.Conn(), query, taskId, target.ResolvedPath, payload)
	if res.IsFailure {
		return "", apperror.WrapSimple(res.Error, "cmdrun.enqueue_audit: insert queue")
	}

	return taskId, nil
}

// CompleteRunTaskAudit records task completion into TaskHistory and updates TaskQueue status.
func CompleteRunTaskAudit(taskId, targetPath, payload string, durationMs int64) error {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "cmdrun.complete_audit: open tasks db")
	}
	defer tasksDB.Close()

	forwardPayload := payload
	if forwardPayload == "" {
		forwardPayload = fmt.Sprintf("duration_ms=%d", durationMs)
	}

	insertErr := tasksDB.InsertTaskHistory(taskId, "run", "file-run", targetPath, forwardPayload, "", "completed")
	if insertErr != nil {
		return apperror.WrapSimple(insertErr, "cmdrun.complete_audit: insert history")
	}

	updateQuery := `UPDATE TaskQueue SET Status = 'completed', UpdatedAt = CURRENT_TIMESTAMP WHERE QueueId = ?`
	_ = store.ExecWrapper(tasksDB.Conn(), updateQuery, taskId)

	return nil
}

// QueryRunTaskHistory returns recent task audit records from TaskHistory.
func QueryRunTaskHistory(limit int) ([]model.TaskHistoryRecord, error) {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdrun.query_history: open tasks db")
	}
	defer tasksDB.Close()

	return tasksDB.ListTaskHistory("run", limit, 0)
}
