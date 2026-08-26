// Package cmd — os_cmd.go provides root CLI entrypoints for OS maintenance commands.
package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Manage operating system updates, upgrades, and repository mirrors",
}

var osUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update OS packages and security patches",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteOSUpdate(context.Background())
	},
}

var osFullUpgradeCmd = &cobra.Command{
	Use:   "full-upgrade",
	Short: "Execute a full OS distribution version upgrade",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteOSFullUpgrade(context.Background())
	},
}

var osFixMirrorsCmd = &cobra.Command{
	Use:     "fix-mirrors",
	Aliases: []string{"update-fix", "fix-update"},
	Short:   "Fix regional repository mirror glitches by switching to canonical US mirrors",
	RunE: func(cmd *cobra.Command, args []string) error {
		return FixRegionalMirrors("")
	},
}

func init() {
	osCmd.AddCommand(osUpdateCmd)
	osCmd.AddCommand(osFullUpgradeCmd)
	osCmd.AddCommand(osFixMirrorsCmd)
}
