package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var fallbackHistCounter uint64

func resolveHistoryID(id string, hostIP string) string {
	if id != "" {
		return id
	}

	seq := atomic.AddUint64(&fallbackHistCounter, 1)

	return fmt.Sprintf("hist-%s-%d-%d", hostIP, time.Now().UnixNano(), seq)
}

func recordSSHHistory(ctx context.Context, execer sqlContextExecer, history SSHHistory) error {
	history.ID = resolveHistoryID(history.ID, history.HostIP)
	history.JoinedAt = resolveSSHHostCreatedAt(history.JoinedAt)
	query := `INSERT INTO ssh_history (id, host_ip, joined_at, user) VALUES (?, ?, ?, ?)`
	_, err := execer.ExecContext(ctx, query, history.ID, history.HostIP, history.JoinedAt, history.User)
	if err != nil {
		return apperror.Wrap(err, "recordSSHHistory", map[string]any{"id": history.ID})
	}

	return nil
}

func commitEnrollTx(tx *sql.Tx, hostID string) error {
	if err := tx.Commit(); err != nil {
		return apperror.Wrap(err, "EnrollSSHHost_Commit", map[string]any{"hostId": hostID})
	}

	return nil
}

func runEnrollTransaction(ctx context.Context, tx *sql.Tx, host SSHHost, history SSHHistory) error {
	if err := UpsertSSHHost(ctx, host, tx); err != nil {
		_ = tx.Rollback()

		return apperror.Wrap(err, "EnrollSSHHost_Upsert", map[string]any{"hostId": host.ID})
	}

	if err := recordSSHHistory(ctx, tx, history); err != nil {
		_ = tx.Rollback()

		return apperror.Wrap(err, "EnrollSSHHost_History", map[string]any{"historyId": history.ID})
	}

	return commitEnrollTx(tx, host.ID)
}

// EnrollSSHHost wraps host upsert and history logging into an atomic transaction.
func EnrollSSHHost(ctx context.Context, host SSHHost, history SSHHistory, db *sql.DB) error {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap(err, "EnrollSSHHost_BeginTx", map[string]any{"hostId": host.ID})
	}

	return runEnrollTransaction(ctx, tx, host, history)
}

// LogSSHHistory logs an SSH connection into the database defensively.
func LogSSHHistory(ctx context.Context, h SSHHistory, db *sql.DB) error {
	ctx = safeContext(ctx)
	_ = EnsureSSHTables(db)

	return LogSSHJoin(ctx, h, db)
}
