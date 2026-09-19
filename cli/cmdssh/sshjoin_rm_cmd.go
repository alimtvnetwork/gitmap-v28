package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

const msgMissingRmTarget = `missing machine alias, IP, or --all to remove

Usage:
  gitmap ssh rm <alias|ip|--all>
  gitmap sj rm <alias|ip|--all>

Examples:
  gitmap ssh rm devbox
  gitmap ssh rm 192.168.1.14
  gitmap ssh rm --all`

func resolveRmTarget(args []string) string {
	for _, arg := range args {
		if isAllFlag(arg) {
			return "all"
		}
	}

	for _, arg := range args {
		clean := strings.TrimSpace(arg)
		if clean != "" && !strings.HasPrefix(clean, "-") {
			return clean
		}
	}

	return ""
}

func validateRmTarget(args []string) (string, error) {
	target := resolveRmTarget(args)
	if target == "" {
		return "", apperror.NewValidationError(msgMissingRmTarget)
	}

	return target, nil
}

func purgeAllSSHRecords(ctx context.Context, sqlDB *sql.DB) error {
	_, errHosts := store.DeleteAllSSHHosts(ctx, sqlDB)
	if errHosts != nil {
		return apperror.WrapSimple(errHosts, "purgeAllSSHRecords.DeleteAllSSHHosts")
	}

	errConn := db.DeleteAllSSHConnections(ctx, sqlDB)
	if errConn != nil {
		return apperror.WrapSimple(errConn, "purgeAllSSHRecords.DeleteAllSSHConnections")
	}

	return nil
}

func purgeTargetSSHRecords(ctx context.Context, sqlDB *sql.DB, target string) error {
	_, errHosts := store.DeleteHostByAliasOrIP(ctx, target, sqlDB)
	if errHosts != nil {
		return apperror.WrapSimple(errHosts, "purgeTargetSSHRecords.DeleteHostByAliasOrIP")
	}

	errConn := db.DeleteSSHConnectionByTarget(ctx, sqlDB, target)
	if errConn != nil {
		return apperror.WrapSimple(errConn, "purgeTargetSSHRecords.DeleteSSHConnectionByTarget")
	}

	return nil
}

func deleteHostRecord(ctx context.Context, target string) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("runSJRm", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()

	if isAllFlag(target) {
		return purgeAllSSHRecords(ctx, dbConn.SQL())
	}

	return purgeTargetSSHRecords(ctx, dbConn.SQL(), target)
}

func printRmSuccess(target string) {
	if isAllFlag(target) {
		fmt.Println("✓ All machines removed from SSH registry.")
		return
	}

	fmt.Printf("✓ Machine '%s' removed from SSH registry.\n", target)
}

// runSJRm handles removing an SSH host by alias or IP, or all hosts with --all.
//
//nolint:revive
func runSJRm(cmd *cobra.Command, args []string, ctx context.Context) error {
	target, err := validateRmTarget(args)
	if err != nil {
		return err
	}

	if err := deleteHostRecord(ctx, target); err != nil {
		return err
	}

	printRmSuccess(target)
	return nil
}

func executeSJRm(ctx context.Context, target string, isForced bool) error {
	_ = isForced
	return runSJRm(nil, []string{target}, ctx)
}
