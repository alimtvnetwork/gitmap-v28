// Package cmd — installer.go defines the root installer CLI command.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

// installerCmd represents the root command for managing installation scripts.
var installerCmd = &cobra.Command{
	Use:   "installer",
	Short: "Manage, create, export, import, and version installation scripts",
	Long:  "A script installation creation, export, import, and versioning system.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstaller(cmd, args)
	},
}

// runInstaller executes the root installer command logic.
func runInstaller(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstaller", "E_INSTALLER_NIL_COMMAND", map[string]any{
			"error": "command is nil",
		})
	}
	if err := cmd.Help(); err != nil {
		appErr := apperror.Wrap(err, "runInstaller", map[string]any{
			"args": args,
		})
		appErr.Code = "E_INSTALLER_COMMAND_FAILED"
		return appErr
	}
	return nil
}
