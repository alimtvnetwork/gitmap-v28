package cmdssh

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

const sqlCreateSSHHistory = `CREATE TABLE IF NOT EXISTS ssh_task_history (
	task_id TEXT PRIMARY KEY,
	action TEXT NOT NULL,
	target TEXT NOT NULL,
	payload_json TEXT NOT NULL,
	created_at TEXT NOT NULL,
	restored_at TEXT DEFAULT ''
)`

func initSSHHistorySchema(conn *sql.DB) error {
	errCfg := store.ConfigureSQLiteConn(conn)
	hasCfgErr := errCfg != nil
	if hasCfgErr {
		return apperror.WrapSimple(errCfg, "openSSHHistoryDB.configure")
	}

	_, errExec := conn.Exec(sqlCreateSSHHistory)
	hasExecErr := errExec != nil
	if hasExecErr {
		return apperror.WrapSimple(errExec, "openSSHHistoryDB.createSchema")
	}

	return nil
}

func openSSHHistoryDB() (*sql.DB, error) {
	dbPath := store.ResolveSplitDbPath("history", "task", "")
	conn, err := sql.Open("sqlite", dbPath)
	hasErr := err != nil
	if hasErr {
		return nil, apperror.WrapSimple(err, "openSSHHistoryDB")
	}

	errInit := initSSHHistorySchema(conn)
	hasInitErr := errInit != nil
	if hasInitErr {
		_ = conn.Close()
		return nil, errInit
	}

	return conn, nil
}

func SnapshotSSHNodesBeforeRemoval(ctx context.Context, action, target string, nodes []store.SSHHost) (string, error) {
	hasNodes := len(nodes) > 0
	if !hasNodes {
		return "", nil
	}

	db, err := openSSHHistoryDB()
	hasDbErr := err != nil
	if hasDbErr {
		return "", err
	}
	defer db.Close()

	return insertHistorySnapshot(ctx, db, action, target, nodes)
}

func insertHistorySnapshot(ctx context.Context, db *sql.DB, action, target string, nodes []store.SSHHost) (string, error) {
	payload, errMarshal := json.Marshal(nodes)
	hasMarshalErr := errMarshal != nil
	if hasMarshalErr {
		return "", apperror.WrapSimple(errMarshal, "insertHistorySnapshot.Marshal")
	}

	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	createdAt := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO ssh_task_history (task_id, action, target, payload_json, created_at) VALUES (?, ?, ?, ?, ?)`
	_, errExec := db.ExecContext(ctx, query, taskID, action, target, string(payload), createdAt)
	hasExecErr := errExec != nil
	if hasExecErr {
		return "", apperror.WrapSimple(errExec, "insertHistorySnapshot.Exec")
	}

	return taskID, nil
}
