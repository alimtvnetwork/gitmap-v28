package cmdssh

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func scanHistoryTask(row *sql.Row) (*SSHHistoryTask, error) {
	var task SSHHistoryTask
	var payload, created string
	err := row.Scan(&task.TaskID, &task.Action, &task.Target, &payload, &created, &task.RestoredAt)
	hasErr := err != nil
	if hasErr {
		return nil, err
	}

	var nodes []store.SSHHost
	_ = json.Unmarshal([]byte(payload), &nodes)
	task.Nodes = nodes
	task.CreatedAt, _ = time.Parse(time.RFC3339, created)

	return &task, nil
}

func GetLastSSHHistoryTask(ctx context.Context) (*SSHHistoryTask, error) {
	db, err := openSSHHistoryDB()
	hasErr := err != nil
	if hasErr {
		return nil, err
	}
	defer db.Close()

	query := `SELECT task_id, action, target, payload_json, created_at, restored_at FROM ssh_task_history WHERE restored_at = '' ORDER BY created_at DESC LIMIT 1`
	row := db.QueryRowContext(ctx, query)

	return scanHistoryTask(row)
}

func GetSSHHistoryTaskByID(ctx context.Context, taskID string) (*SSHHistoryTask, error) {
	db, err := openSSHHistoryDB()
	hasErr := err != nil
	if hasErr {
		return nil, err
	}
	defer db.Close()

	query := `SELECT task_id, action, target, payload_json, created_at, restored_at FROM ssh_task_history WHERE task_id = ? LIMIT 1`
	row := db.QueryRowContext(ctx, query, taskID)

	return scanHistoryTask(row)
}

func MarkSSHHistoryTaskRestored(ctx context.Context, taskID string) error {
	db, err := openSSHHistoryDB()
	hasErr := err != nil
	if hasErr {
		return err
	}
	defer db.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE ssh_task_history SET restored_at = ? WHERE task_id = ?`
	_, errExec := db.ExecContext(ctx, query, now, taskID)
	hasExecErr := errExec != nil
	if hasExecErr {
		return apperror.WrapSimple(errExec, "MarkSSHHistoryTaskRestored")
	}

	return nil
}
