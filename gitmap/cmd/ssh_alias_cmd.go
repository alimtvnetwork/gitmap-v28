package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// SSHAliasCmd represents the gitmap ssh as command.
var SSHAliasCmd = &cobra.Command{
	Use:   "as [ip] [alias name]",
	Short: "Create an SSH alias for an IP using the 'as' keyword",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHAlias(cmd, args, cmd.Context())
	},
}

//nolint:revive
func runSSHAlias(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) < 2 {
		return apperror.NewSimple("runSSHAlias", "E_INTERNAL_ERROR")
	}

	ip := args[0]
	aliasName := args[1]

	return saveAliasCommand(ctx, ip, aliasName)
}

func saveAliasCommand(ctx context.Context, ip string, alias string) error {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, "saveAliasCommand")
	}
	defer db.Close()

	wrapper, appErr := dbengine.WrapDb(db.Conn(), dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	return runSaveAliasTx(ctx, wrapper, ip, alias)
}

func runSaveAliasTx(ctx context.Context, wrapper *dbengine.DbWrapper, ip, alias string) error {
	txErr := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeSaveAliasTx(ctx, tx, ip, alias)
	})
	if txErr != nil {
		return txErr
	}

	fmt.Printf("Successfully saved alias %s for %s\n", alias, ip)

	return nil
}

func executeSaveAliasTx(ctx context.Context, tx *dbengine.TxWrapper, ip, alias string) *apperror.AppError {
	host := store.SSHHost{
		ID:        fmt.Sprintf("ssh-%d", time.Now().UnixNano()),
		Alias:     alias,
		IP:        ip,
		Username:  "",
		CreatedAt: time.Now().UTC(),
	}
	if err := store.InsertSSHHost(ctx, host, tx.Tx()); err != nil {
		return apperror.WrapSimple(err, "saveAliasCommand")
	}

	return nil
}
