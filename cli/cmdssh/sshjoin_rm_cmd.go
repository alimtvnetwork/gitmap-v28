package cmdssh

import (
	"context"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

const msgMissingRmTarget = `missing machine alias or IP to remove

Usage:
  gitmap ssh join rm <alias|ip>
  gitmap sj rm <alias|ip>

Examples:
  gitmap sj rm devbox
  gitmap sj rm 192.168.1.14`

func validateRmTarget(args []string) (string, error) {
	if len(args) == 0 {
		return "", apperror.NewValidationError(msgMissingRmTarget)
	}

	target := strings.TrimSpace(args[0])
	if target == "" {
		return "", apperror.NewValidationError(msgMissingRmTarget)
	}

	return target, nil
}

func deleteHostRecord(ctx context.Context, target string) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("runSJRm", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()

	_, err = store.DeleteHostByAliasOrIP(ctx, target, dbConn.SQL())
	return err
}

// runSJRm handles removing an SSH host by alias or IP.
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

	fmt.Printf("✓ Machine '%s' removed from SSH registry.\n", target)
	return nil
}

func executeSJRm(ctx context.Context, target string, isForced bool) error {
	return runSJRm(nil, []string{target}, ctx)
}
