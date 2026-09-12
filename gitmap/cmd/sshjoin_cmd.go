package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func executeSSHJoin(ctx context.Context, target string, history store.SSHHistory) error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to open db", "err": err.Error()})
	}
	defer dbConn.Close()

	return runJoinTransaction(ctx, dbConn.SQL(), target, history)
}

func runJoinTransaction(ctx context.Context, db *sql.DB, target string, history store.SSHHistory) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "begin tx error", "err": err.Error()})
	}
	defer tx.Rollback()

	if err := insertJoinRecords(ctx, tx, target, history); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "commit tx error", "err": err.Error()})
	}

	fmt.Println("Joined successfully")

	return nil
}

func insertJoinRecords(ctx context.Context, tx *sql.Tx, target string, history store.SSHHistory) error {
	host := store.SSHHost{
		ID:       history.ID,
		IP:       history.HostIP,
		Username: history.User,
		Alias:    target,
	}

	if err := store.InsertSSHHost(ctx, host, tx); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "insert host error", "err": err.Error()})
	}

	return logSSHJoinInTx(ctx, tx, history)
}

func logSSHJoinInTx(ctx context.Context, tx *sql.Tx, history store.SSHHistory) error {
	query := `INSERT INTO ssh_history (id, host_ip, joined_at, user) VALUES (?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query, history.ID, history.HostIP, history.JoinedAt, history.User)
	if err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "log join error", "err": err.Error()})
	}

	return nil
}
