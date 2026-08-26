// Package cmd — installer_os_update.go defines dedicated per-OS installer update CLI commands.
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/installer"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/spf13/cobra"
)

func createOSUpdateCmd(osTarget string) *cobra.Command {
	return &cobra.Command{
		Use:                fmt.Sprintf("update-%s <slug>", osTarget),
		Short:              fmt.Sprintf("Update installer script record specifically for %s", osTarget),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeOSUpdate(args, osTarget)
		},
	}
}

func executeOSUpdate(args []string, osTarget string) error {
	if len(args) == 0 {
		return apperror.New("executeOSUpdate", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug is required",
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

	if _, err := mgr.Update(slug, osTarget); err != nil {
		return err
	}

	fmt.Printf("Installer \"%s\" (%s) updated successfully.\n", slug, osTarget)
	return nil
}

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(createOSUpdateCmd(constants.OSTargetUbuntu))
		installerCmd.AddCommand(createOSUpdateCmd(constants.OSTargetDebian))
		installerCmd.AddCommand(createOSUpdateCmd(constants.OSTargetArch))
		installerCmd.AddCommand(createOSUpdateCmd(constants.OSTargetCentOS))
		installerCmd.AddCommand(createOSUpdateCmd(constants.OSTargetFedora))
		installerCmd.AddCommand(createOSUpdateCmd(constants.OSTargetMac))
		installerCmd.AddCommand(createOSUpdateCmd(constants.OSTargetUnix))
	}
}
