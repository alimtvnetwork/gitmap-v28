// Package cmd — installer_os_cmds.go defines dedicated per-OS installer CLI commands.
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/installer"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/spf13/cobra"
)

func createOSInstallCmd(osTarget string) *cobra.Command {
	return &cobra.Command{
		Use:                fmt.Sprintf("install-%s <slug>", osTarget),
		Short:              fmt.Sprintf("Execute installer script targeted for %s", osTarget),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeOSInstall(args, osTarget)
		},
	}
}

func executeOSInstall(args []string, osTarget string) error {
	if len(args) == 0 {
		return apperror.New("executeOSInstall", "E_INSTALLER_INVALID_INPUT", map[string]any{
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

	ctx := context.Background()
	if err := mgr.ExecuteOrdered(ctx, slug, osTarget); err != nil {
		return err
	}

	fmt.Printf("Installer \"%s\" (%s) executed successfully.\n", slug, osTarget)
	return nil
}

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetUbuntu))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetDebian))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetArch))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetCentOS))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetFedora))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetMac))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetUnix))
	}
}
