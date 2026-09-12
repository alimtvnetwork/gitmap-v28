package cmdssh

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
	wrapper, appErr := dbengine.WrapDb(db, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	txErr := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return insertJoinRecords(ctx, tx, target, history)
	})
	if txErr != nil {
		return txErr
	}

	fmt.Println("Joined successfully")

	return nil
}

func insertJoinRecords(ctx context.Context, tx *dbengine.TxWrapper, target string, history store.SSHHistory) *apperror.AppError {
	host := store.SSHHost{
		ID:       history.ID,
		IP:       history.HostIP,
		Username: history.User,
		Alias:    target,
	}

	if err := store.InsertSSHHost(ctx, host, tx.Tx()); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "insert host error", "err": err.Error()})
	}

	return logSSHJoinInTx(ctx, tx, history)
}

func logSSHJoinInTx(ctx context.Context, tx *dbengine.TxWrapper, history store.SSHHistory) *apperror.AppError {
	query := `INSERT INTO ssh_history (id, host_ip, joined_at, user) VALUES (?, ?, ?, ?)`
	if _, err := tx.Exec(ctx, query, history.ID, history.HostIP, history.JoinedAt, history.User); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "log join error", "err": err.Error()})
	}

	return nil
}

var SSHJoinCmd = &cobra.Command{
	Use:     "ssh-join",
	Aliases: []string{"sj", "ssh-joined", "ssh-joiner"},
	Short:   "Join an SSH machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHJoin(cmd, args, cmd.Context())
	},
}

//nolint:revive
func runSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) > 0 {
		switch args[0] {
		case "add", "rm", "ls", "history", "add-auth":
			// Handled by subcommands
		default:
			return apperror.New("runSSHJoin", "E_INTERNAL_ERROR", map[string]any{"arg": args[0]})
		}
	}

	return nil
}

var SJRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove machine-alias or ip",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJRm(cmd, args, cmd.Context())
	},
}

//nolint:revive
func runSJRm(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) != 1 {
		return apperror.New("runSJRm", "E_INTERNAL_ERROR", map[string]any{"msg": "invalid argument count"})
	}

	return nil
}

func init() {
	SSHJoinCmd.AddCommand(SJRmCmd)
	SSHJoinCmd.AddCommand(SJAddAuthCmd)
	SSHJoinCmd.AddCommand(SJLsCmd)
	SSHJoinCmd.AddCommand(SJHistCmd)
}
