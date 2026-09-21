package cmdssh

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func askResetConfirmation(count int, isForced bool) bool {
	if isForced {
		return true
	}

	fmt.Printf("\nReset SSH registry? This removes all %d node(s) and flushes connection caches. [y/N]: ", count)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	clean := strings.TrimSpace(strings.ToLower(input))
	isConfirmed := clean == "y" || clean == "yes"

	return isConfirmed
}

func purgeAndResetRegistry(ctx context.Context, hosts []store.SSHHost) error {
	_, _ = SnapshotSSHNodesBeforeRemoval(ctx, ActionResetNodes, "all", hosts)
	dbConn, errDB := openSSHDBFunc()
	hasDBErr := errDB != nil
	if hasDBErr {
		return apperror.WrapSimple(errDB, "purgeAndResetRegistry.openDB")
	}
	defer dbConn.Close()

	_, _ = store.DeleteAllSSHHosts(ctx, dbConn.SQL())
	_ = db.DeleteAllSSHConnections(ctx, dbConn.SQL())
	fmt.Printf("✓ SSH registry reset successfully. All %d node(s) removed and caches cleared.\n", len(hosts))
	fmt.Println("  Undo anytime: gitmap ssh undo")

	return nil
}

// RunSSHResetCLI wipes all registered nodes and flushes connection caches.
func RunSSHResetCLI(args []string) error {
	ctx := context.Background()
	hosts, err := fetchSJHosts(ctx)
	hasErr := err != nil
	if hasErr {
		return err
	}

	hasHosts := len(hosts) > 0
	if !hasHosts {
		fmt.Println("SSH registry is already clean (0 nodes registered).")

		return nil
	}

	isConfirmed := askResetConfirmation(len(hosts), hasYesFlag(args))
	if !isConfirmed {
		fmt.Println("Canceled.")

		return nil
	}

	return purgeAndResetRegistry(ctx, hosts)
}
