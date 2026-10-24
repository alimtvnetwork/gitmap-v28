package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

var openSSHDBFunc = openDB

const msgMissingJoinTarget = `missing target host

Usage:
  gitmap ssh-join <user@ip|ip> [alias] [flags]
  gitmap ssh-join add <user@ip|ip> [alias] [flags]
  gitmap ssh-join add-with-pass <user@ip|ip> [password] [alias] [flags]
  gitmap ssh join <user@ip|ip> [alias] [flags]
  gitmap sj <user@ip|ip> [alias] [flags]

Examples:
  › 1. Join Machine by User and IP (user@ip)
  # Join user@ip with alias (prompts interactively for password on first join and stores in salted RSA vault; leave blank to skip saving)
  gitmap ssh join alim@192.168.1.14 devbox
  gitmap ssh-join user@192.168.1.14
  gitmap ssh-join root@192.168.1.14 prod-server
  gitmap ssh-join add dev@192.168.1.50 devbox
  gitmap ssh-join add-with-pass alim@192.168.1.14 secret123 devbox
  gitmap ssh-join 192.168.1.14
  gitmap ssh join administrator@192.168.1.20 main    # prompts for password & saves in salted RSA vault
  gitmap ssh join ubuntu@192.168.1.14:2222 prod --auth`

func executeSSHJoin(ctx context.Context, target string, history store.SSHHistory) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to open db", "err": err.Error()})
	}

	defer dbConn.Close()

	return runJoinTransaction(ctx, dbConn.SQL(), target, history)
}

func executeJoinInTx(ctx context.Context, wrapper *dbengine.DBWrapper, target string, history store.SSHHistory) error {
	appErr := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return insertJoinRecords(ctx, tx, target, history)
	})
	if appErr != nil {
		return appErr
	}

	return nil
}

func runJoinTransaction(ctx context.Context, db *sql.DB, target string, history store.SSHHistory) error {
	wrapper, appErr := dbengine.WrapDB(db, dbengine.DBSQLite)
	if appErr != nil {
		return appErr
	}

	if txErr := executeJoinInTx(ctx, wrapper, target, history); txErr != nil {
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
	if execRes := tx.Exec(ctx, query, history.ID, history.HostIP, history.JoinedAt, history.User); execRes.IsFailure() {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "log join error", "err": execRes.AppError().Error()})
	}

	return nil
}

func isSJClusterImportSubcommand(sub string) bool {
	return sub == "import-cluster" || sub == "import" || sub == "import-config"
}

func isSJAddWithPassSubcommand(sub string) bool {
	return sub == "add-with-pass" || sub == "add-pass" || sub == "add-password"
}

func isSJAddSubcommand(sub string) bool {
	return sub == "add" || sub == "join" || sub == "new" || sub == "enroll"
}

func isSJBootstrapSubcommand(sub string) bool {
	return sub == "bootstrap" || sub == "bs"
}

func isSJInstallSubcommand(sub string) bool {
	return sub == "install"
}

func routeSJClusterSpecialSubcommand(args []string) (bool, error) {
	if isSJInstallSubcommand(args[0]) {
		return true, RunClusterInstallCLI(args[1:])
	}
	if isSJBootstrapSubcommand(args[0]) {
		return true, RunClusterBootstrapCLI(args[1:])
	}
	if isSJClusterImportSubcommand(args[0]) {
		return true, RunClusterImportCLI(args[1:])
	}
	return false, nil
}

func routeSSHJoinSpecialSubcommand(cmd *cobra.Command, args []string) (bool, error) {
	if isHandled, err := routeSJClusterSpecialSubcommand(args); isHandled {
		return true, err
	}
	if isSJAddWithPassSubcommand(args[0]) {
		return true, executeEnrollWithPassCLI(resolveContext(cmd.Context()), args[1:])
	}
	if isSJAddSubcommand(args[0]) {
		return true, executeEnrollCLI(resolveContext(cmd.Context()), args[1:])
	}
	return false, nil
}

func routeSSHJoinCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return RunSSHJoinCLI(args)
	}
	if isHandled, err := routeSSHJoinSpecialSubcommand(cmd, args); isHandled {
		return err
	}
	if !isSJSubcommand(args[0]) {
		return executeEnrollCLI(resolveContext(cmd.Context()), args)
	}
	return RunSSHJoinCLI(args)
}

var SSHJoinCmd = &cobra.Command{
	Use:     "ssh-join [user@ip|ip] [alias]",
	Aliases: []string{"sj", "ssh-joined", "ssh-joiner"},
	Short:   "Join an SSH machine by user@ip or IP address",
	Args:    cobra.ArbitraryArgs,
	RunE:    routeSSHJoinCmd,
}

//nolint:revive
func runSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error {
	return RunSSHJoinCLI(args)
}

func isSJScanSubcommand(sub string) bool {
	return sub == "scan" || sub == "find" || sub == "discover" || sub == "probe"
}

func isSJStatusSubcommand(sub string) bool {
	return sub == "status" || sub == "ping" || sub == "health" || sub == "check"
}

func isSJSubcommand(sub string) bool {
	if isSJAddSubcommand(sub) || isSJAddWithPassSubcommand(sub) || isSJScanSubcommand(sub) || isSJStatusSubcommand(sub) || isSJClusterImportSubcommand(sub) || isSJBootstrapSubcommand(sub) || isSJInstallSubcommand(sub) {
		return true
	}

	return sub == "ls" || sub == "list" || sub == "nodes" || sub == "node" || sub == "rm" || sub == "remove" || sub == "delete" ||
		sub == "clear" || sub == "reset" ||
		sub == "add-auth" || sub == "auth" || sub == "history" || sub == "hist" ||
		sub == "auth-key" || sub == "copy-id" || sub == "fix-auth"
}

func isSJListSubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "nodes" || sub == "node"
}

func executeSJList(ctx context.Context) error {
	return printSJList(ctx, os.Stdout, 0)
}

func dispatchSJBasicSubcommand(ctx context.Context, sub string, args []string) (bool, error) {
	if isSJInstallSubcommand(sub) {
		return true, RunClusterInstallCLI(args)
	}
	if isSJClusterImportSubcommand(sub) {
		return true, RunClusterImportCLI(args)
	}
	if isSJAddWithPassSubcommand(sub) {
		return true, executeEnrollWithPassCLI(ctx, args)
	}
	if isSJAddSubcommand(sub) {
		return true, executeEnrollCLI(ctx, args)
	}
	return false, nil
}

func dispatchSJQuerySubcommand(ctx context.Context, sub string, args []string) (bool, error) {
	if isSJScanSubcommand(sub) {
		return true, RunSJScan(SJScanCmd, args, ctx)
	}
	if isSJStatusSubcommand(sub) {
		return true, RunSJStatus(SJStatusCmd, args, ctx)
	}
	if isSJListSubcommand(sub) {
		return true, executeSJList(ctx)
	}
	return false, nil
}

func dispatchSJRemovalSubcommand(ctx context.Context, sub string, args []string) (bool, error) {
	if sub == "rm" || sub == "remove" || sub == "delete" {
		return true, runSJRm(nil, args, ctx)
	}
	if sub == "clear" {
		return true, runSSHClearNodes(args)
	}
	if sub == "reset" {
		return true, RunSSHResetCLI(args)
	}

	return false, nil
}

func dispatchSJActionSubcommand(ctx context.Context, sub string, args []string) (bool, error) {
	if isHandled, err := dispatchSJRemovalSubcommand(ctx, sub, args); isHandled {
		return true, err
	}
	if sub == "add-auth" || sub == "auth" {
		return true, runSJAddAuth(nil, args, ctx)
	}
	if sub == "auth-key" || sub == "copy-id" || sub == "fix-auth" {
		return true, RunSSHAuthKeyDeployCLI(args)
	}
	if sub == "pass" || sub == "password" {
		return true, RunSSHPassCLI(args)
	}

	return false, nil
}

func dispatchSJSubcommand(ctx context.Context, sub string, args []string) error {
	if isHandled, err := dispatchSJBasicSubcommand(ctx, sub, args); isHandled {
		return err
	}
	if isHandled, err := dispatchSJQuerySubcommand(ctx, sub, args); isHandled {
		return err
	}
	if isHandled, err := dispatchSJActionSubcommand(ctx, sub, args); isHandled {
		return err
	}
	return runSJHistory(nil, args, ctx)
}
