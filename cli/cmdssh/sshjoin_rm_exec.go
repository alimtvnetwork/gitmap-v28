package cmdssh

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func purgeCandidateNodes(ctx context.Context, sqlDB *sql.DB, candidates []store.SSHHost) error {
	for _, cand := range candidates {
		_, _ = store.DeleteHostByAliasOrIP(ctx, cand.Alias, sqlDB)
		_, _ = store.DeleteHostByAliasOrIP(ctx, cand.IP, sqlDB)
		_ = db.DeleteSSHConnectionByTarget(ctx, sqlDB, cand.Alias)
		_ = db.DeleteSSHConnectionByTarget(ctx, sqlDB, cand.IP)
	}

	return nil
}

func executeNodeDeletion(ctx context.Context, action, target string, candidates []store.SSHHost) error {
	taskID, errSnap := SnapshotSSHNodesBeforeRemoval(ctx, action, target, candidates)
	hasSnapErr := errSnap != nil
	if hasSnapErr {
		return errSnap
	}

	dbConn, errDB := openSSHDBFunc()
	hasDBErr := errDB != nil
	if hasDBErr {
		return apperror.WrapSimple(errDB, "executeNodeDeletion.openDB")
	}
	defer dbConn.Close()

	_ = purgeCandidateNodes(ctx, dbConn.SQL(), candidates)
	printNodeDeletionSuccess(len(candidates), taskID)

	return nil
}

func printNodeDeletionSuccess(count int, taskID string) {
	_ = taskID
	fmt.Printf("✓ %d machine(s) removed from SSH registry.\n", count)
	fmt.Println("  Undo anytime: gitmap ssh undo")
}

func runSSHClearNodes(args []string) error {
	ctx := context.Background()
	hosts, err := fetchSJHosts(ctx)
	hasErr := err != nil
	if hasErr {
		return err
	}

	hasHosts := len(hosts) > 0
	if !hasHosts {
		fmt.Println("No registered SSH nodes found.")

		return nil
	}

	isConfirmed := askRemovalConfirmation(hosts, hasYesFlag(args))
	if !isConfirmed {
		fmt.Println("Canceled.")

		return nil
	}

	return executeNodeDeletion(ctx, ActionClearNodes, "all", hosts)
}

func dispatchRmNodesArgs(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs {
		return runSSHClearNodes(args)
	}

	target := args[0]
	isAll := target == "clear" || target == "all"
	if isAll {
		return runSSHClearNodes(args[1:])
	}

	return executeTargetNodeRm(target, args[1:])
}

func executeSJRm(ctx context.Context, target string, isForced bool) error {
	flags := []string{}
	if isForced {
		flags = append(flags, "-y")
	}

	return executeTargetNodeRm(target, flags)
}
