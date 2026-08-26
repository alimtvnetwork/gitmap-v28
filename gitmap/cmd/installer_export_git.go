// Package cmd — installer_export_git.go defines Git direct export CLI commands.
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/installer"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/spf13/cobra"
)

var installerExportGitCmd = &cobra.Command{
	Use:                "export-git <slug> <folder-or-url>",
	Short:              "Export single installer directly to a local Git directory or remote Git repo",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeExportGit(args, false)
	},
}

var installerExportAllGitCmd = &cobra.Command{
	Use:                "export-all-git <folder-or-url>",
	Aliases:            []string{"export-all -git", "export -all-git"},
	Short:              "Export all installer scripts directly to a local Git directory or remote Git repo",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeExportGit(args, true)
	},
}

func executeExportGit(args []string, isAll bool) error {
	if (isAll && len(args) < 1) || (!isAll && len(args) < 2) {
		return apperror.New("executeExportGit", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "insufficient arguments provided for git export",
		})
	}

	slug := ""
	target := ""
	if isAll {
		target = strings.TrimSpace(args[0])
	} else {
		slug = strings.TrimSpace(args[0])
		target = strings.TrimSpace(args[1])
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		return errMigrate
	}

	mgr, errMgr := installer.NewManager(db)
	if errMgr != nil {
		return errMgr
	}

	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") || strings.HasPrefix(target, "git@") {
		if err := mgr.ExportToRemoteGitRepo(slug, target, "main", "", "", false); err != nil {
			return err
		}
	} else {
		if err := mgr.ExportToGitFolder(slug, target, "", ""); err != nil {
			return err
		}
	}

	fmt.Printf("Exported installer(s) to Git target %q successfully.\n", target)
	return nil
}

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerExportGitCmd)
		installerCmd.AddCommand(installerExportAllGitCmd)
	}
}
