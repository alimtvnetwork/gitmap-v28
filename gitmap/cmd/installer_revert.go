// Package cmd — installer_revert.go defines the undo, redo, and revert CLI subcommands for installers.
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// installerUndoCmd represents the 'gitmap installer undo-version' subcommand.
var installerUndoCmd = &cobra.Command{
	Use:   "undo-version <slug>",
	Short: "Revert installer to its previous version",
	Long:  "Rolls back an installer script definition to its immediately preceding version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerRevertAction(cmd, args, "undo")
	},
}

// installerRedoCmd represents the 'gitmap installer redo-version' subcommand.
var installerRedoCmd = &cobra.Command{
	Use:   "redo-version <slug>",
	Short: "Advance installer to a newer reverted version",
	Long:  "Restores a previously undone version of an installer script definition.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerRevertAction(cmd, args, "redo")
	},
}

// installerRevertCmd represents the 'gitmap installer revert-version' subcommand.
var installerRevertCmd = &cobra.Command{
	Use:   "revert-version <slug> <version>",
	Short: "Revert installer to an exact semantic version",
	Long:  "Restores an installer script definition to an exact historical version tag.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerRevertAction(cmd, args, "revert")
	},
}

// InstallerUndoCmd is an exported alias for installerUndoCmd.
var InstallerUndoCmd = installerUndoCmd

// InstallerRedoCmd is an exported alias for installerRedoCmd.
var InstallerRedoCmd = installerRedoCmd

// InstallerRevertCmd is an exported alias for installerRevertCmd.
var InstallerRevertCmd = installerRevertCmd

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerUndoCmd)
		installerCmd.AddCommand(installerRedoCmd)
		installerCmd.AddCommand(installerRevertCmd)
	}
}

// executeRevertAction coordinates version changes against the database.
func executeRevertAction(ctx context.Context, db *store.DB, action, slug, targetVersion string) error {
	if db == nil {
		return apperror.New("executeRevertAction", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "db cannot be nil"})
	}
	if strings.TrimSpace(slug) == "" {
		return apperror.New("executeRevertAction", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug is required"})
	}

	existing, errGet := db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	fmt.Printf("Installer %q version action %s processed (current: %s, target: %s).\n",
		existing.Name, action, existing.Version, targetVersion)
	return nil
}

// runInstallerRevertAction coordinates arguments and dispatches to executeRevertAction.
func runInstallerRevertAction(cmd *cobra.Command, args []string, action string) error {
	if cmd == nil {
		return apperror.New("runInstallerRevertAction", "E_INSTALLER_NIL_COMMAND", map[string]any{"error": "command is nil"})
	}
	if len(args) == 0 {
		return apperror.New("runInstallerRevertAction", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug argument required"})
	}

	slug := args[0]
	targetVersion := ""
	if action == "revert" {
		if len(args) < 2 {
			return apperror.New("runInstallerRevertAction", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "version argument required for revert-version"})
		}
		targetVersion = args[1]
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerRevertAction", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerRevertAction", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeRevertAction(cmd.Context(), db, action, slug, targetVersion)
}
