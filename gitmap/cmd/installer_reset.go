// Package cmd — installer_reset.go defines the installer reset CLI subcommand.
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

// ResetInstallerFlags encapsulates parsed CLI options for installer resets.
type ResetInstallerFlags struct {
	Slug     string
	ResetAll bool
}

// installerResetCmd represents the 'gitmap reset installer' or 'gitmap installer reset' subcommand.
var installerResetCmd = &cobra.Command{
	Use:   "reset [slug]",
	Short: "Reset installer script records",
	Long:  "Resets installer definitions and version histories by slug or with --all.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerReset(cmd, args)
	},
}

// InstallerResetCmd is an exported alias for installerResetCmd.
var InstallerResetCmd = installerResetCmd

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerResetCmd)
	}
}

// parseInstallerResetFlags parses command-line flags and positional args for installer reset.
func parseInstallerResetFlags(args []string) (*ResetInstallerFlags, error) {
	fs := flag.NewFlagSet("reset", flag.ContinueOnError)
	resetAll := fs.Bool("all", false, "Reset all installers")
	fs.BoolVar(resetAll, "a", false, "Reset all shorthand")

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
		appErr := apperror.Wrap(err, "parseInstallerResetFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"
		return nil, appErr
	}

	slug := ""
	if len(positional) > 0 {
		slug = positional[0]
	}

	if !*resetAll && strings.TrimSpace(slug) == "" {
		return nil, apperror.New("parseInstallerResetFlags", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "slug is required when --all is not specified",
		})
	}

	return &ResetInstallerFlags{
		Slug:     strings.TrimSpace(slug),
		ResetAll: *resetAll,
	}, nil
}

// executeInstallerReset coordinates resetting through the store DB.
func executeInstallerReset(ctx context.Context, db *store.DB, flags *ResetInstallerFlags) error {
	if db == nil {
		return apperror.New("executeInstallerReset", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "db cannot be nil"})
	}

	if errReset := db.ResetInstallers(flags.Slug, flags.ResetAll); errReset != nil {
		appErr := apperror.Wrap(errReset, "executeInstallerReset", map[string]any{"slug": flags.Slug, "all": flags.ResetAll})
		appErr.Code = "E_INSTALLER_RESET_FAILED"
		return appErr
	}

	if flags.ResetAll {
		fmt.Println("All installer script records reset successfully.")
	} else {
		fmt.Printf("Installer %q reset successfully.\n", flags.Slug)
	}
	return nil
}

// runInstallerReset coordinates flag parsing, database connection, and reset execution.
func runInstallerReset(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerReset", "E_INSTALLER_NIL_COMMAND", map[string]any{"error": "command is nil"})
	}

	flags, errParse := parseInstallerResetFlags(args)
	if errParse != nil {
		return errParse
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerReset", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerReset", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeInstallerReset(cmd.Context(), db, flags)
}
