// Package cmd — sshjoin_distribute_cmd.go defines CLI commands for multi-host key distribution.
package cmd

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

var sjDistributeKeysCmd = &cobra.Command{
	Use:     "distribute-keys <ip1,ip2...>",
	Aliases: []string{"broadcast-keys"},
	Short:   "Distribute local SSH public key to multiple remote machines",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeDistributeKeys(args)
	},
}

func executeDistributeKeys(args []string) error {
	if len(args) == 0 {
		return apperror.New("executeDistributeKeys", "E_INVALID_ARGS", map[string]any{"error": "hosts argument required"})
	}

	hosts := ParseMultiIPList(args[0])
	return DistributeKeysToHosts(context.Background(), hosts, "root", 22)
}

func init() {
	if SSHJoinCmd != nil {
		SSHJoinCmd.AddCommand(sjDistributeKeysCmd)
	}
}
