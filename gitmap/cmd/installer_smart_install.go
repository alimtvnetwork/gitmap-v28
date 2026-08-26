package cmd

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/installer"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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
	default:
		return constants.OSTargetUbuntu
	}
}

func validateInstallSlug(args []string) (string, error) {
	if len(args) == 0 {
		return "", apperror.New("validateInstallSlug", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug is required",
		})
	}
	return strings.TrimSpace(args[0]), nil
}

func runSmartInstaller(slug, osTarget string) error {
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
	return mgr.ExecuteOrdered(context.Background(), slug, osTarget)
}

func hasProfilePrefix(slug string) bool {
	_, ok := resolveProfileTree(slug)
	return ok
}

func renderSmartInstallSummary(slug, osTarget string) {
	fmt.Printf("Installer \"%s\" auto-installed for %s successfully.\n", slug, osTarget)
	if hasProfilePrefix(slug) {
		printProfileInstallSummary(slug)
		return
	}
	printInstallSummaryHeader(slug)
}

func executeSmartInstall(args []string) error {
	slug, errValidate := validateInstallSlug(args)
	if errValidate != nil {
		return errValidate
	}
	osTarget := detectCurrentHostOSTarget()
	if errRun := runSmartInstaller(slug, osTarget); errRun != nil {
		return errRun
	}
	renderSmartInstallSummary(slug, osTarget)
	return nil
}

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerSmartInstallCmd)
	}
}
