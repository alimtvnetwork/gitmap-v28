package cmdssh

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// EnqueueSSHTask wraps an SSH operation into the task queue with forward and inverse payloads.
func EnqueueSSHTask(ctx context.Context, action, target, forwardPayload, inversePayload string) (string, error) {
	db, err := openSSHHistoryDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC3339)
	if resQueue := insertSSHTaskQueue(db, taskID, action, target, forwardPayload, inversePayload, now); resQueue != nil {
		return "", resQueue
	}
	if resHist := insertSSHHistoryEntry(db, taskID, action, target, forwardPayload, inversePayload, now); resHist != nil {
		return "", resHist
	}

	recordSSHTaskToRootTasksDB(taskID, action, target, forwardPayload, inversePayload)
	return taskID, nil
}

func insertSSHTaskQueue(db *sql.DB, taskID, action, target, forward, inverse, now string) error {
	queryQueue := `INSERT INTO TaskQueue (QueueId, Section, Action, Target, ForwardPayload, InversePayload, Status, CreatedAt, UpdatedAt) VALUES (?, 'ssh', ?, ?, ?, ?, 'pending', ?, ?)`
	resQueue := store.ExecWrapper(db, queryQueue, taskID, action, target, forward, inverse, now, now)
	if resQueue.IsFailure {
		return apperror.WrapSimple(resQueue.Error, "EnqueueSSHTask.queue")
	}
	return nil
}

func insertSSHHistoryEntry(db *sql.DB, taskID, action, target, forward, inverse, now string) error {
	queryHist := `INSERT INTO TaskHistory (TaskId, Section, Action, Target, ForwardPayload, InversePayload, CreatedAt) VALUES (?, 'ssh', ?, ?, ?, ?, ?)`
	resHist := store.ExecWrapper(db, queryHist, taskID, action, target, forward, inverse, now)
	if resHist.IsFailure {
		return apperror.WrapSimple(resHist.Error, "EnqueueSSHTask.hist")
	}
	return nil
}

// RecordSSHAddNodeTask records an add node operation with inverse delete payload.
func RecordSSHAddNodeTask(ctx context.Context, host store.SSHHost) (string, error) {
	forward, _ := json.Marshal(host)
	inverse, _ := json.Marshal(map[string]string{"alias": host.Alias, "ip": host.IP})
	return EnqueueSSHTask(ctx, ActionAddNode, host.Alias, string(forward), string(inverse))
}

// RecordSSHDeleteKeyTask records an SSH key delete operation with inverse restore payload.
func RecordSSHDeleteKeyTask(ctx context.Context, key model.SSHKey) (string, error) {
	forward, _ := json.Marshal(map[string]string{"name": key.Name})
	inverse, _ := json.Marshal(key)
	return EnqueueSSHTask(ctx, ActionRmKey, key.Name, string(forward), string(inverse))
}
