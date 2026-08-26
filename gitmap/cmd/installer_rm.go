// Package cmd — installer_rm.go provides deletion and version purging for installers.
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/installer"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

var installerRmCmd = &cobra.Command{
	Use:                "rm <slug>",
	Aliases:            []string{"delete"},
	Short:              "Delete an installer script record or specific OS target",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeInstallerRm(args)
	},
}

func executeInstallerRm(args []string) error {
	if len(args) == 0 {
		return apperror.New("executeInstallerRm", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "installer slug is required",
		})
	}

	slug := strings.TrimSpace(args[0])
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

	if err := mgr.Delete(slug); err != nil {
		return err
	}

	fmt.Printf("Installer \"%s\" deleted successfully.\n", slug)
	return nil
}

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerRmCmd)
	}
}
