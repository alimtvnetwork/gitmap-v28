package cmd

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func executeSJRm(ctx context.Context, target string, force bool) error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return apperror.New("executeSJRm", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to open db", "err": err.Error()})
	}
	defer dbConn.Close()

	if err := store.DeleteHostByIP(ctx, target, dbConn.SQL()); err != nil {
		return apperror.New("executeSJRm", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to delete host", "err": err.Error()})
	}

	fmt.Println("Machine removed")
	return nil
}
