package cmdssh

import (
	"context"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var SJRmCmd = &cobra.Command{
	Use:   "rm [nodes|keys|alias|ip|all] [flags]",
	Short: "Remove registered SSH nodes or keys with confirmation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJRm(cmd, args, cmd.Context())
	},
}

func dispatchRmCategory(ctx context.Context, target string, args []string) (bool, error) {
	isKeys := target == "keys" || target == "key"
	if isKeys {
		return true, runSSHDelete(args[1:])
	}

	isClear := target == "clear" || target == "all"
	if isClear {
		return true, runSSHClearNodes(args[1:])
	}

	isNodes := target == "nodes" || target == "node"
	if isNodes {
		return true, dispatchRmNodesArgs(args[1:])
	}

	return false, nil
}

func executeTargetNodeRm(target string, flags []string) error {
	ctx := context.Background()
	hosts, err := fetchSJHosts(ctx)
	hasErr := err != nil
	if hasErr {
		return err
	}

	candidates := filterCandidateNodes(target, hosts)
	hasNone := len(candidates) == 0
	if hasNone {
		return apperror.NewValidationError(fmt.Sprintf("no SSH node matching '%s' found. Use 'gitmap ssh nodes' to view registered nodes.", target))
	}

	isConfirmed := askRemovalConfirmation(candidates, hasYesFlag(flags))
	if !isConfirmed {
		fmt.Println("Canceled.")

		return nil
	}

	return executeNodeDeletion(ctx, ActionRmNode, target, candidates)
}

func runSJRm(cmd *cobra.Command, args []string, ctx context.Context) error {
	_ = cmd
	hasArgs := len(args) > 0
	if !hasArgs {
		return promptInteractiveRM(ctx, args)
	}

	target, flags := parseRmTokens(args)
	isHandled, err := dispatchRmCategory(ctx, target, args)
	if isHandled {
		return err
	}

	return executeTargetNodeRm(target, flags)
}
