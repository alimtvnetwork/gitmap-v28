package store

import (
	"context"
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// LogSSHJoin logs an SSH connection into the database.
//nolint:revive
func LogSSHJoin(ctx context.Context, h SSHHistory, db *sql.DB) error {
	query := `INSERT INTO ssh_history (id, host_ip, joined_at, user) VALUES (?, ?, ?, ?)`
	_, err := db.ExecContext(ctx, query, h.ID, h.HostIP, h.JoinedAt, h.User)
	if err != nil {
		return &apperror.AppError{
			Op:    "LogSSHJoin",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"id": h.ID, "host_ip": h.HostIP},
		}
	}
	return nil
}

// ListSSHHistory retrieves a list of SSH history records, ordered by joined_at descending.
//nolint:revive
func ListSSHHistory(ctx context.Context, limit int, offset int, db *sql.DB) ([]SSHHistory, error) {
	query := `SELECT id, host_ip, joined_at, user FROM ssh_history ORDER BY joined_at DESC LIMIT ? OFFSET ?`
	rows, err := db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, &apperror.AppError{
			Op:    "ListSSHHistory",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"limit": limit, "offset": offset},
		}
	}
	defer rows.Close()

	var result []SSHHistory
	for rows.Next() {
		var h SSHHistory
		if err := rows.Scan(&h.ID, &h.HostIP, &h.JoinedAt, &h.User); err != nil {
			return nil, &apperror.AppError{
				Op:    "ListSSHHistory_Scan",
				Code:  "E_INTERNAL_ERROR",
				Cause: err,
				Ctx:   nil,
			}
		}
		result = append(result, h)
	}

	if err := rows.Err(); err != nil {
		return nil, &apperror.AppError{
			Op:    "ListSSHHistory_Rows",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   nil,
		}
	}

	if result == nil {
		result = []SSHHistory{}
	}

	return result, nil
}
