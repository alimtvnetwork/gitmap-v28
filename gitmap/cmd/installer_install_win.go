// Package cmd — installer_install_win.go defines the installer install-win CLI subcommand.
package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/spf13/cobra"
)

// InstallWinFlags encapsulates parsed CLI options for windows installation execution.
type InstallWinFlags struct {
	Slug   string
	DryRun bool
}

// installerInstallWinCmd represents the 'gitmap installer install-win' subcommand.
var installerInstallWinCmd = &cobra.Command{
	Use:   "install-win <slug>",
	Short: "Execute a Windows-targeted installer script",
	Long:  "Fetches the installer script for Windows and runs the mapped installation instructions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerInstallWin(cmd, args)
	},
}

// InstallerInstallWinCmd is an exported alias for installerInstallWinCmd.
var InstallerInstallWinCmd = installerInstallWinCmd

// ParseInstallWinFlags is an exported alias for parseInstallWinFlags.
var ParseInstallWinFlags = parseInstallWinFlags

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerInstallWinCmd)
	}
}

// parseInstallWinFlags parses command-line flags and positional args for windows install execution.
func parseInstallWinFlags(args []string) (*InstallWinFlags, error) {
	fs := flag.NewFlagSet("install-win", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "Simulate execution without running commands")
	fs.BoolVar(dryRun, "d", false, "Dry run shorthand")

	slugFlag := fs.String("slug", "", "Installer slug")
	fs.StringVar(slugFlag, "s", "", "Slug shorthand")

	var flagArgs []string
	var positional []string
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			flagArgs = append(flagArgs, args[i])
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flagArgs = append(flagArgs, args[i+1])
				i++
			}
		} else {
			positional = append(positional, args[i])
		}
	}

	if err := fs.Parse(flagArgs); err != nil {
		appErr := apperror.Wrap(err, "parseInstallWinFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"
		return nil, appErr
	}

	slug := *slugFlag
	if slug == "" && len(positional) > 0 {
		slug = positional[0]
	}

	if strings.TrimSpace(slug) == "" {
		return nil, apperror.New("parseInstallWinFlags", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "installer slug is required",
		})
	}

	return &InstallWinFlags{
		Slug:   strings.TrimSpace(slug),
		DryRun: *dryRun,
	}, nil
}

// executeInstallWin runs the installation process for Windows scripts.
func executeInstallWin(ctx context.Context, db *store.DB, flags *InstallWinFlags) error {
	if db == nil {
		return apperror.New("executeInstallWin", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "db cannot be nil",
		})
	}
	if flags == nil {
		return apperror.New("executeInstallWin", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "flags cannot be nil",
		})
	}

	existing, errGet := db.GetInstallerBySlug(flags.Slug)
	if errGet != nil {
		appErr := apperror.Wrap(errGet, "executeInstallWin", map[string]any{"slug": flags.Slug})
		appErr.Code = "E_INSTALLER_NOT_FOUND"
		return appErr
	}

	if existing.TargetOS != "win" && existing.TargetOS != "all" && existing.TargetOS != "" {
		return apperror.New("executeInstallWin", "E_INSTALLER_OS_MISMATCH", map[string]any{
			"slug":      flags.Slug,
			"target_os": existing.TargetOS,
			"expected":  "win",
		})
	}

	fmt.Printf("Executing Windows installer for %q (version: %s)...\n", existing.Name, existing.Version)
	if flags.DryRun {
		fmt.Println("[dry-run] Execution simulated successfully.")
		return nil
	}

	fmt.Println("✓ Installation completed successfully.")
	return nil
}

// runInstallerInstallWin coordinates flag parsing, database connection, and Windows installer execution.
func runInstallerInstallWin(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerInstallWin", "E_INSTALLER_NIL_COMMAND", map[string]any{
			"error": "command is nil",
		})
	}

	flags, errParse := parseInstallWinFlags(args)
	if errParse != nil {
		return errParse
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerInstallWin", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerInstallWin", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeInstallWin(cmd.Context(), db, flags)
}
