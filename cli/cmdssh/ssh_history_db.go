package cmdssh

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

func migrateLegacySSHHistoryFile(targetPath string) {
	legacyPath := store.ResolveSplitDbPath("history", "task", "")
	fi, errStat := os.Stat(legacyPath)
	isLegacyPresent := errStat == nil && !fi.IsDir()
	if isLegacyPresent {
		copyLegacyFile(legacyPath, targetPath)
	}
}

func copyLegacyFile(src, dst string) {
	s, errOpen := os.Open(src)
	if errOpen != nil {
		return
	}
	defer s.Close()

	d, errCreate := os.Create(dst)
	if errCreate != nil {
		return
	}
	defer d.Close()

	_, _ = io.Copy(d, s)
}

func resolveSSHHistoryDBPath() string {
	targetPath := store.ResolveSectionTasksDbPath("ssh", "")
	fi, errStat := os.Stat(targetPath)
	isExisting := errStat == nil && !fi.IsDir()
	if !isExisting {
		_ = os.MkdirAll(filepath.Dir(targetPath), 0755)
		migrateLegacySSHHistoryFile(targetPath)
	}
	return targetPath
}

func openSSHHistoryDB() (*sql.DB, error) {
	_ = resolveSSHHistoryDBPath()
	secDB, err := store.OpenSectionTasksDB("ssh", "")
	if err != nil {
		return nil, apperror.WrapSimple(err, "openSSHHistoryDB.OpenSectionTasksDB")
	}

	conn := secDB.Conn()
	migrateLegacySSHHistoryTables(conn)

	return conn, nil
}

func migrateLegacySSHHistoryTables(conn *sql.DB) {
	if checkTableExists(conn, "ssh_task_history") {
		_ = store.ExecWrapper(conn, `INSERT OR IGNORE INTO TaskHistory 
			(TaskId, Section, Action, Target, ForwardPayload, InversePayload, RestoredAt, CreatedAt)
			SELECT task_id, 'ssh', action, target, forward_payload, inverse_payload, COALESCE(restored_at, ''), created_at 
			FROM ssh_task_history;`)
	}
	if checkTableExists(conn, "ssh_task_queue") {
		_ = store.ExecWrapper(conn, `INSERT OR IGNORE INTO TaskQueue 
			(QueueId, Section, Action, Target, ForwardPayload, InversePayload, Status, CreatedAt, UpdatedAt)
			SELECT task_id, 'ssh', action, target, forward_payload, inverse_payload, status, created_at, updated_at 
			FROM ssh_task_queue;`)
	}
}

func checkTableExists(conn *sql.DB, tableName string) bool {
	query := `SELECT 1 FROM sqlite_master WHERE type='table' AND name=? LIMIT 1`
	row := conn.QueryRow(query, tableName)
	var dummy int
	return row.Scan(&dummy) == nil
}

// SnapshotSSHNodesBeforeRemoval stores nodes prior to deletion for undo support.
func SnapshotSSHNodesBeforeRemoval(ctx context.Context, action, target string, nodes []store.SSHHost) (string, error) {
	hasNodes := len(nodes) > 0
	if !hasNodes {
		return "", nil
	}

	db, err := openSSHHistoryDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	return insertHistorySnapshot(ctx, db, action, target, nodes)
}

func insertHistorySnapshot(ctx context.Context, db *sql.DB, action, target string, nodes []store.SSHHost) (string, error) {
	payload, errMarshal := json.Marshal(nodes)
	if errMarshal != nil {
		return "", apperror.WrapSimple(errMarshal, "insertHistorySnapshot.Marshal")
	}

	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	createdAt := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO TaskHistory (TaskId, Section, Action, Target, ForwardPayload, InversePayload, CreatedAt) VALUES (?, 'ssh', ?, ?, ?, ?, ?)`
	res := store.ExecWrapper(db, query, taskID, action, target, target, string(payload), createdAt)
	if res.IsFailure {
		return "", apperror.WrapSimple(res.Error, "insertHistorySnapshot.Exec")
	}

	recordSSHTaskToRootTasksDB(taskID, action, target, target, string(payload))
	return taskID, nil
}

func recordSSHTaskToRootTasksDB(taskID, action, target, fwd, inv string) {
	rdb, errOpen := store.OpenTasksRootSplitDB()
	if errOpen != nil {
		return
	}
	defer rdb.Close()

	sqlQuery := `INSERT OR IGNORE INTO TaskHistory 
		(TaskId, Section, Action, Target, ForwardPayload, InversePayload, Status) 
		VALUES (?, 'ssh', ?, ?, ?, ?, 'completed')`
	res := store.ExecWrapper(rdb.Conn(), sqlQuery, taskID, action, target, fwd, inv)
	_, _ = res.Destruct() // lint-allow: ignore-db-error reason="optional cross-split sync"
}

// EnqueueSSHTask wraps an SSH operation into the task queue with forward and inverse payloads.
func EnqueueSSHTask(ctx context.Context, action, target, forwardPayload, inversePayload string) (string, error) {
	db, err := openSSHHistoryDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC3339)
	resQueue := insertSSHTaskQueue(db, taskID, action, target, forwardPayload, inversePayload, now)
	if resQueue != nil {
		return "", resQueue
	}

	resHist := insertSSHHistoryEntry(db, taskID, action, target, forwardPayload, inversePayload, now)
	if resHist != nil {
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
