package cmdssh

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func hasYesFlag(args []string) bool {
	for _, a := range args {
		isYes := a == "-y" || a == "--yes" || a == "-f" || a == "--force"
		if isYes {
			return true
		}
	}

	return false
}

func askRemovalConfirmation(candidates []store.SSHHost, isForced bool) bool {
	if isForced {
		return true
	}

	fmt.Println("\nCandidate node(s) for removal:")
	_ = RenderSSHHostsTable(os.Stdout, candidates)
	fmt.Printf("\nAre you sure you want to remove %d node(s)? [y/N]: ", len(candidates))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	clean := strings.TrimSpace(strings.ToLower(input))
	isConfirmed := clean == "y" || clean == "yes"

	return isConfirmed
}

func promptInteractiveRM(ctx context.Context, args []string) error {
	fmt.Println("\nWhat would you like to remove?")
	fmt.Println("  [1] Registered SSH nodes")
	fmt.Println("  [2] Managed SSH keys")
	fmt.Print("\nSelect option (1-2) or 'q' to quit: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	return dispatchInteractiveRM(ctx, choice, args)
}

func dispatchInteractiveRM(ctx context.Context, choice string, args []string) error {
	isNodes := choice == "1" || choice == "nodes"
	if isNodes {
		return runSSHClearNodes(args)
	}

	isKeys := choice == "2" || choice == "keys"
	if isKeys {
		return runSSHDelete(args)
	}

	fmt.Println("Canceled.")

	return nil
}
