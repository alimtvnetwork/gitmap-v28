// Package cmdpipeline — pipeline_tasks_db.go manages pipeline section task recordings.
package cmdpipeline

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RecordPipelineCheckInTask logs a pipeline error check-in task into pipeline-tasks.db and tasks root DB.
func RecordPipelineCheckInTask(repo string, payload PipelineErrorLogsPayload) *apperror.AppError {
	fwdJSON, invJSON := serializePipelinePayload(payload)
	taskId := generatePipelineTaskId(repo)
	action := "pipeline_error_checkin"

	errSection := recordToPipelineSectionDB(taskId, action, repo, fwdJSON, invJSON)
	if errSection != nil {
		return errSection
	}

	return recordToTasksRootDB(taskId, action, repo, fwdJSON, invJSON)
}

func serializePipelinePayload(payload PipelineErrorLogsPayload) (string, string) {
	bytes, err := json.Marshal(payload)
	hasErr := err != nil
	if hasErr {
		return "{}", "{}"
	}

	return string(bytes), "{}"
}

func generatePipelineTaskId(repo string) string {
	nowNano := time.Now().UnixNano()

	return fmt.Sprintf("pipe_%s_%d", repo, nowNano)
}

func recordToPipelineSectionDB(taskId, action, target, fwd, inv string) *apperror.AppError {
	pdb, errOpen := store.OpenSectionTasksDB("pipeline", "")
	hasOpenErr := errOpen != nil
	if hasOpenErr {
		return apperror.WrapSimple(errOpen, "open pipeline tasks db")
	}

	defer pdb.Close()

	return insertTaskHistoryRecord(pdb.Conn(), taskId, "pipeline", action, target, fwd, inv)
}

func recordToTasksRootDB(taskId, action, target, fwd, inv string) *apperror.AppError {
	rdb, errOpen := store.OpenTasksRootSplitDB()
	hasOpenErr := errOpen != nil
	if hasOpenErr {
		return apperror.WrapSimple(errOpen, "open tasks root db")
	}

	defer rdb.Close()

	return insertTaskHistoryRecord(rdb.Conn(), taskId, "pipeline", action, target, fwd, inv)
}

func insertTaskHistoryRecord(conn *sql.DB, taskId, section, action, target, fwd, inv string) *apperror.AppError {
	sqlQuery := `INSERT OR IGNORE INTO TaskHistory
		(TaskId, Section, Action, Target, ForwardPayload, InversePayload, Status)
		VALUES (?, ?, ?, ?, ?, ?, 'completed')`
	res := store.ExecWrapper(conn, sqlQuery, taskId, section, action, target, fwd, inv)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "insert task history")
	}

	return nil
}
