package cmd

import (
	"context"
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

	tx, err := dbConn.SQL().BeginTx(ctx, nil)
	if err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "begin tx error", "err": err.Error()})
	}
	defer tx.Rollback()

	host := store.SSHHost{
		ID:       history.ID,
		IP:       history.HostIP,
		Username: history.User,
		Alias:    target,
	}

	if err := store.InsertSSHHost(ctx, host, tx); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "insert host error", "err": err.Error()})
	}

	if err := tx.Commit(); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "commit tx error", "err": err.Error()})
	}

	if err := store.LogSSHJoin(ctx, history, dbConn.SQL()); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "log join error", "err": err.Error()})
	}

	fmt.Println("Joined successfully")
	return nil
}
