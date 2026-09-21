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

const (
	sqlCreateSSHHistory = `CREATE TABLE IF NOT EXISTS ssh_task_history (
	task_id TEXT PRIMARY KEY,
	action TEXT NOT NULL,
	target TEXT NOT NULL,
	payload_json TEXT NOT NULL,
	forward_payload TEXT NOT NULL DEFAULT '',
	inverse_payload TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL,
	restored_at TEXT DEFAULT ''
);`

	sqlCreateSSHTaskQueue = `CREATE TABLE IF NOT EXISTS ssh_task_queue (
	task_id TEXT PRIMARY KEY,
	action TEXT NOT NULL,
	target TEXT NOT NULL,
	forward_payload TEXT NOT NULL DEFAULT '',
	inverse_payload TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'pending',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);`
)

func initSSHHistorySchema(conn *sql.DB) error {
	errCfg := store.ConfigureSQLiteConn(conn)
	if errCfg != nil {
		return apperror.WrapSimple(errCfg, "openSSHHistoryDB.configure")
	}

	_, _ = conn.Exec(sqlCreateSSHHistory)
	_, _ = conn.Exec(sqlCreateSSHTaskQueue)
	_, _ = conn.Exec(`ALTER TABLE ssh_task_history ADD COLUMN forward_payload TEXT NOT NULL DEFAULT ''`)
	_, _ = conn.Exec(`ALTER TABLE ssh_task_history ADD COLUMN inverse_payload TEXT NOT NULL DEFAULT ''`)

	return nil
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

func migrateLegacySSHHistoryFile(targetPath string) {
	legacyPath := store.ResolveSplitDbPath("history", "task", "")
	fi, errStat := os.Stat(legacyPath)
	isLegacyPresent := errStat == nil && !fi.IsDir()
	if isLegacyPresent {
		copyLegacyFile(legacyPath, targetPath)
	}
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
	dbPath := resolveSSHHistoryDBPath()
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "openSSHHistoryDB")
	}

	errInit := initSSHHistorySchema(conn)
	if errInit != nil {
		_ = conn.Close()
		return nil, errInit
	}

	return conn, nil
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
	query := `INSERT INTO ssh_task_history (task_id, action, target, payload_json, forward_payload, inverse_payload, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, errExec := db.ExecContext(ctx, query, taskID, action, target, string(payload), target, string(payload), createdAt)
	if errExec != nil {
		return "", apperror.WrapSimple(errExec, "insertHistorySnapshot.Exec")
	}

	return taskID, nil
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
	queryQueue := `INSERT INTO ssh_task_queue (task_id, action, target, forward_payload, inverse_payload, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 'pending', ?, ?)`
	_, _ = db.ExecContext(ctx, queryQueue, taskID, action, target, forwardPayload, inversePayload, now, now)

	queryHist := `INSERT INTO ssh_task_history (task_id, action, target, payload_json, forward_payload, inverse_payload, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, errHist := db.ExecContext(ctx, queryHist, taskID, action, target, inversePayload, forwardPayload, inversePayload, now)
	if errHist != nil {
		return "", apperror.WrapSimple(errHist, "EnqueueSSHTask.hist")
	}

	return taskID, nil
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
