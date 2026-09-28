// Package cmdos — os_cmd.go provides root CLI entrypoints for OS maintenance, machine identity, and alias commands.
package cmdos

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Manage operating system updates, machine identity, network alias, and system hygiene",
	RunE: func(cmd *cobra.Command, args []string) error {
		RenderModernOSHelp()
		return nil
	},
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

var osMachineCmd = &cobra.Command{
	Use:                "machine",
	Aliases:            []string{"machines", "machine-name", "hostname"},
	Short:              "View, set, change, or revert OS machine name and SSH fleet hostnames",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunMachineCLI(args)
	},
}

var osAliasCmd = &cobra.Command{
	Use:                "alias",
	Aliases:            []string{"aliases"},
	Short:              "View, set, change, or revert machine network alias across local and SSH fleet",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAliasCLI(args)
	},
}

func init() {
	osCmd.AddCommand(osUpdateCmd)
	osCmd.AddCommand(osFullUpgradeCmd)
	osCmd.AddCommand(osFixMirrorsCmd)
	osCmd.AddCommand(osMachineCmd)
	osCmd.AddCommand(osAliasCmd)
	osCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		RenderModernOSHelp()
	})
}

// RunOSCLI routes OS CLI commands including machine, alias, modernized help, and system tools.
func RunOSCLI(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) {
		RenderModernOSHelp()
		return nil
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	if isHandled, err := tryDispatchOSIdentityOrCobra(sub, args); isHandled {
		return err
	}
	return runOS(args)
}

func tryDispatchOSIdentityOrCobra(sub string, args []string) (bool, error) {
	switch sub {
	case "machine", "machines", "machine-name", "hostname":
		return true, RunMachineCLI(args[1:])
	case "alias", "aliases":
		return true, RunAliasCLI(args[1:])
	case "info", "os-info", "osinfo", "sysinfo", "system-info":
		return true, RunOSInfoCLI(args[1:])
	case "update", "full-upgrade", "fix-mirrors", "update-fix", "fix-update":
		osCmd.SetArgs(args)
		return true, osCmd.Execute()
	default:
		return false, nil
	}
}
