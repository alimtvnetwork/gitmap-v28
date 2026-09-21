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
	var payload, fwd, inv, created string
	err := row.Scan(&task.TaskID, &task.Action, &task.Target, &payload, &fwd, &inv, &created, &task.RestoredAt)
	if err != nil {
		return nil, err
	}

	var nodes []store.SSHHost
	_ = json.Unmarshal([]byte(payload), &nodes)
	task.Nodes = nodes
	task.ForwardPayload = fwd
	task.InversePayload = inv
	task.CreatedAt, _ = time.Parse(time.RFC3339, created)

	return &task, nil
}

// GetLastSSHHistoryTask retrieves the most recent non-restored SSH history task.
func GetLastSSHHistoryTask(ctx context.Context) (*SSHHistoryTask, error) {
	db, err := openSSHHistoryDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT task_id, action, target, payload_json, forward_payload, inverse_payload, created_at, restored_at FROM ssh_task_history WHERE restored_at = '' ORDER BY created_at DESC LIMIT 1`
	row := db.QueryRowContext(ctx, query)

	return scanHistoryTask(row)
}

// GetLastRestoredSSHHistoryTask retrieves the most recent restored SSH history task for redo.
func GetLastRestoredSSHHistoryTask(ctx context.Context) (*SSHHistoryTask, error) {
	db, err := openSSHHistoryDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT task_id, action, target, payload_json, forward_payload, inverse_payload, created_at, restored_at FROM ssh_task_history WHERE restored_at != '' ORDER BY restored_at DESC LIMIT 1`
	row := db.QueryRowContext(ctx, query)

	return scanHistoryTask(row)
}

// GetSSHHistoryTaskByID retrieves an SSH history task by its unique ID.
func GetSSHHistoryTaskByID(ctx context.Context, taskID string) (*SSHHistoryTask, error) {
	db, err := openSSHHistoryDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT task_id, action, target, payload_json, forward_payload, inverse_payload, created_at, restored_at FROM ssh_task_history WHERE task_id = ? LIMIT 1`
	row := db.QueryRowContext(ctx, query, taskID)

	return scanHistoryTask(row)
}

// MarkSSHHistoryTaskRestored sets the restored_at timestamp for a task.
func MarkSSHHistoryTaskRestored(ctx context.Context, taskID string) error {
	db, err := openSSHHistoryDB()
	if err != nil {
		return err
	}
	defer db.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE ssh_task_history SET restored_at = ? WHERE task_id = ?`
	_, errExec := db.ExecContext(ctx, query, now, taskID)
	if errExec != nil {
		return apperror.WrapSimple(errExec, "MarkSSHHistoryTaskRestored")
	}

	return nil
}

// MarkSSHHistoryTaskActive clears the restored_at timestamp for redo.
func MarkSSHHistoryTaskActive(ctx context.Context, taskID string) error {
	db, err := openSSHHistoryDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE ssh_task_history SET restored_at = '' WHERE task_id = ?`
	_, errExec := db.ExecContext(ctx, query, taskID)
	if errExec != nil {
		return apperror.WrapSimple(errExec, "MarkSSHHistoryTaskActive")
	}

	return nil
}
