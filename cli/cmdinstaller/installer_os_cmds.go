// Package cmd — installer_os_cmds.go defines dedicated per-OS installer CLI commands.
package cmdinstaller

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/installer"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func createOSInstallCmd(osTarget string) *cobra.Command {
	return &cobra.Command{
		Use:                fmt.Sprintf("install-%s <slug> [--force-all]", osTarget),
		Short:              fmt.Sprintf("Execute installer script targeted for %s", osTarget),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeOSInstall(args, osTarget)
		},
	}
}

func executeOSInstall(args []string, osTarget string) error {
	slug, hasForceAll, errInput := parseInstallArgs(args)
	if errInput != nil {
		return errInput
	}

	isCompat, skipMsg := isHostCompatibleWithTarget(osTarget, hasForceAll)
	if !isCompat {
		fmt.Printf("[gitmap] %s\n", skipMsg)
		return nil
	}

	return runManagerInstall(slug, osTarget)
}

func parseInstallArgs(args []string) (string, bool, error) {
	if len(args) == 0 {
		return "", false, apperror.New("executeOSInstall", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug is required",
		})
	}
	hasForce := false
	var slug string
	for _, a := range args {
		if a == "--force-all" || a == "-f" || a == "force-all" {
			hasForce = true
		} else if slug == "" {
			slug = strings.TrimSpace(a)
		}
	}
	if slug == "" {
		return "", false, apperror.New("executeOSInstall", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug is required",
		})
	}
	return slug, hasForce, nil
}

func isHostCompatibleWithTarget(target string, hasForceAll bool) (bool, string) {
	if hasForceAll {
		return true, ""
	}
	isWinHost := runtime.GOOS == constants.OSWindows
	if isWinHost && isUnixTarget(target) {
		return false, fmt.Sprintf("Skipping Unix installer target %q on Windows (use --force-all to override)", target)
	}
	if !isWinHost && target == constants.OSTargetWin {
		return false, "Skipping Windows installer target on POSIX/Unix (use --force-all to override)"
	}
	return true, ""
}

func isUnixTarget(target string) bool {
	switch target {
	case constants.OSTargetUbuntu,
		constants.OSTargetDebian,
		constants.OSTargetCentOS,
		constants.OSTargetFedora,
		constants.OSTargetArch,
		constants.OSTargetUnix,
		constants.OSTargetMac:
		return true
	default:
		return false
	}
}

func runManagerInstall(slug, osTarget string) error {
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
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetWin))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetUbuntu))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetDebian))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetArch))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetCentOS))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetFedora))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetMac))
		installerCmd.AddCommand(createOSInstallCmd(constants.OSTargetUnix))
	}
}
