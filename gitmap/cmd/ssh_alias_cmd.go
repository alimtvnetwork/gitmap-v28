package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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

//nolint:revive
func saveAliasCommand(ctx context.Context, ip string, alias string) error {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, "saveAliasCommand")
	}
	defer db.Close()

	tx, err := db.Conn().BeginTx(ctx, nil)
	if err != nil {
		return apperror.WrapSimple(err, "saveAliasCommand")
	}
	defer tx.Rollback()

	id := fmt.Sprintf("ssh-%d", time.Now().UnixNano())
	host := store.SSHHost{
		ID:        id,
		Alias:     alias,
		IP:        ip,
		Username:  "",
		CreatedAt: time.Now().UTC(),
	}

	if err := store.InsertSSHHost(ctx, host, tx); err != nil {
		return apperror.WrapSimple(err, "saveAliasCommand")
	}

	if err := tx.Commit(); err != nil {
		return apperror.WrapSimple(err, "saveAliasCommand")
	}

	fmt.Printf("Successfully saved alias %s for %s\n", alias, ip)
	return nil
}
