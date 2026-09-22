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
	var fwd, inv, created string
	err := row.Scan(&task.TaskID, &task.Action, &task.Target, &fwd, &inv, &created, &task.RestoredAt)
	if err != nil {
		return nil, err
	}

	var nodes []store.SSHHost
	_ = json.Unmarshal([]byte(inv), &nodes)
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

	query := `SELECT TaskId, Action, Target, ForwardPayload, InversePayload, CreatedAt, RestoredAt
		FROM TaskHistory
		WHERE Section = 'ssh' AND RestoredAt = ''
		ORDER BY CreatedAt DESC LIMIT 1`
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

	query := `SELECT TaskId, Action, Target, ForwardPayload, InversePayload, CreatedAt, RestoredAt
		FROM TaskHistory
		WHERE Section = 'ssh' AND RestoredAt != ''
		ORDER BY RestoredAt DESC LIMIT 1`
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

	query := `SELECT TaskId, Action, Target, ForwardPayload, InversePayload, CreatedAt, RestoredAt
		FROM TaskHistory
		WHERE TaskId = ? LIMIT 1`
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
	query := `UPDATE TaskHistory SET RestoredAt = ? WHERE TaskId = ?`
	res := store.ExecWrapper(db, query, now, taskID)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "MarkSSHHistoryTaskRestored")
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

	query := `UPDATE TaskHistory SET RestoredAt = '' WHERE TaskId = ?`
	res := store.ExecWrapper(db, query, taskID)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "MarkSSHHistoryTaskActive")
	}

	return nil
}

// SSHHistoryRecord represents a summarized history entry for CLI display.
type SSHHistoryRecord struct {
	TaskId     string
	Action     string
	Target     string
	Status     string
	RestoredAt string
	CreatedAt  string
}

// ListSSHHistoryTasks queries SSH tasks from TaskHistory with limit and offset.
func ListSSHHistoryTasks(ctx context.Context, limit, offset int) ([]SSHHistoryRecord, error) {
	db, err := openSSHHistoryDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT TaskId, Action, Target, Status, RestoredAt, CreatedAt
		FROM TaskHistory
		WHERE Section = 'ssh'
		ORDER BY TaskHistoryId DESC
		LIMIT ? OFFSET ?`
	rows, errQuery := db.QueryContext(ctx, query, limit, offset)
	if errQuery != nil {
		return nil, apperror.WrapSimple(errQuery, "ListSSHHistoryTasks.Query")
	}
	defer rows.Close()

	var results []SSHHistoryRecord
	for rows.Next() {
		var rec SSHHistoryRecord
		if errScan := rows.Scan(&rec.TaskId, &rec.Action, &rec.Target, &rec.Status, &rec.RestoredAt, &rec.CreatedAt); errScan == nil {
			results = append(results, rec)
		}
	}

	return results, nil
}
