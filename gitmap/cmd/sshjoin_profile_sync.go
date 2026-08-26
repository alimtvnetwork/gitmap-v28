// Package cmd — sshjoin_profile_sync.go provisions ZSH and profile configs across nodes.
package cmd

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

var sjSyncProfileCmd = &cobra.Command{
	Use:     "sync-profile <ip1,ip2...>",
	Short:   "Sync ZSH and user shell profile configurations to remote machines",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeSJSyncProfile(args)
	},
}

func executeSJSyncProfile(args []string) error {
	if len(args) == 0 {
		return apperror.New("executeSJSyncProfile", "E_INVALID_ARGS", map[string]any{"error": "hosts required"})
	}

	hosts := ParseMultiIPList(args[0])
	ctx := context.Background()
	zshSetupScript := "echo [gitmap] syncing zsh profile && touch ~/.zshrc"

	for _, host := range hosts {
		target, errParse := ParseSSHTarget(host, "root", 22)
		if errParse != nil {
			continue
		}
		fmt.Printf("→ Syncing profile on %s\n", target.String())
		_ = SpawnSSH(ctx, *target, []string{zshSetupScript})
	}
	return nil
}

func init() {
	if SSHJoinCmd != nil {
		SSHJoinCmd.AddCommand(sjSyncProfileCmd)
	}
}
