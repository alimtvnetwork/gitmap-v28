// Package cmd — sshjoin_broadcast.go broadcasts remote shell commands across nodes.
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

var sjBroadcastCmd = &cobra.Command{
	Use:                "broadcast <ip1,ip2...> <command>",
	Short:              "Broadcast a shell command across multiple SSH nodes",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeSJBroadcast(args)
	},
}

func executeSJBroadcast(args []string) error {
	if len(args) < 2 {
		return apperror.New("executeSJBroadcast", "E_INVALID_ARGS", map[string]any{"error": "hosts and command required"})
	}

	hosts := ParseMultiIPList(args[0])
	command := args[1]
	ctx := context.Background()

	for _, host := range hosts {
		target, errParse := ParseSSHTarget(host, "root", 22)
		if errParse != nil {
			fmt.Printf("⚠️ Parse error for %s: %v\n", host, errParse)
			continue
		}
		fmt.Printf("→ Executing on %s: %s\n", target.String(), command)
		_ = SpawnSSH(ctx, *target, []string{command})
	}
	return nil
}

func init() {
	if SSHJoinCmd != nil {
		SSHJoinCmd.AddCommand(sjBroadcastCmd)
	}
}
