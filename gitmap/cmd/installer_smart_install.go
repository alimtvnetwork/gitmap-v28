// Package cmd — installer_smart_install.go provides smart auto-detect installation.
package cmd

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/installer"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/spf13/cobra"
)

var installerSmartInstallCmd = &cobra.Command{
	Use:                "install <slug>",
	Short:              "Smart auto-detect host OS and execute matching installer script",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeSmartInstall(args)
	},
}

func detectCurrentHostOSTarget() string {
	switch runtime.GOOS {
	case "windows":
		return constants.OSTargetWin
	case "darwin":
		return constants.OSTargetMac
	default: // linux / unix
		return constants.OSTargetUbuntu // default linux target or unix
	}
}

func executeSmartInstall(args []string) error {
	if len(args) == 0 {
		return apperror.New("executeSmartInstall", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug is required",
		})
	}

	slug := strings.TrimSpace(args[0])
	osTarget := detectCurrentHostOSTarget()

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

	fmt.Printf("Installer \"%s\" auto-installed for %s successfully.\n", slug, osTarget)
	return nil
}

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerSmartInstallCmd)
	}
}
