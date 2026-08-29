// Package cmd — installer_update.go defines the installer update CLI subcommand.
package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// UpdateInstallerFlags encapsulates parsed CLI options for installer updates.
type UpdateInstallerFlags struct {
	Slug         string
	Description  string
	TargetOS     string
	Version      string
	Instructions string
}

// installerUpdateCmd represents the 'gitmap installer update' subcommand.
var installerUpdateCmd = &cobra.Command{
	Use:                "update <slug>",
	Short:              "Update an existing installer script record",
	Long:               "Updates an installer script record and bumps its version in history.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerUpdate(cmd, args)
	},
}

// InstallerUpdateCmd is an exported alias for installerUpdateCmd.
var InstallerUpdateCmd = installerUpdateCmd

// ParseUpdateFlags is an exported alias for parseUpdateFlags.
var ParseUpdateFlags = parseUpdateFlags

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerUpdateCmd)
	}
}

// parseUpdateFlags parses command-line flags and positional args for installer update.
func parseUpdateFlags(args []string) (*UpdateInstallerFlags, error) {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	desc := fs.String("description", "", "Description for the installer script")
	fs.StringVar(desc, "desc", "", "Description shorthand")
	fs.StringVar(desc, "d", "", "Description shorthand")

	targetOS := fs.String("os", "", "Target operating system")
	fs.StringVar(targetOS, "target-os", "", "Target operating system")
	fs.StringVar(targetOS, "o", "", "Target operating system shorthand")

	ver := fs.String("version", "", "Version for the installer script")
	fs.StringVar(ver, "v", "", "Version shorthand")

	instr := fs.String("instructions", "", "Installation instructions or JSON payload")
	fs.StringVar(instr, "script", "", "Instructions shorthand")
	fs.StringVar(instr, "i", "", "Instructions shorthand")

	slugFlag := fs.String("slug", "", "Installer slug")
	fs.StringVar(slugFlag, "s", "", "Slug shorthand")

	flagArgs, positional := separateFlagAndPositionalArgs(args)

	if err := fs.Parse(flagArgs); err != nil {
		appErr := apperror.Wrap(err, "parseUpdateFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"
		return nil, appErr
	}

	slug := *slugFlag
	if slug == "" && len(positional) > 0 {
		slug = positional[0]
	}

	if strings.TrimSpace(slug) == "" {
		return nil, apperror.New("parseUpdateFlags", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "installer slug is required",
		})
	}

	return &UpdateInstallerFlags{
		Slug:         strings.TrimSpace(slug),
		Description:  strings.TrimSpace(*desc),
		TargetOS:     strings.TrimSpace(*targetOS),
		Version:      strings.TrimSpace(*ver),
		Instructions: strings.TrimSpace(*instr),
	}, nil
}

// executeUpdate executes the update logic against the SQLite database.

func executeInstallerUpdate(ctx context.Context, db *store.DB, flags *UpdateInstallerFlags) error {
	if db == nil {
		return apperror.New("executeUpdate", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "db cannot be nil",
		})
	}
	if flags == nil {
		return apperror.New("executeUpdate", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "flags cannot be nil",
		})
	}

	existing, errGet := db.GetInstallerBySlug(flags.Slug)
	if errGet != nil {
		appErr := apperror.Wrap(errGet, "executeUpdate", map[string]any{"slug": flags.Slug})
		appErr.Code = "E_INSTALLER_NOT_FOUND"
		return appErr
	}

	if flags.Description != "" {
		existing.Description = flags.Description
	}
	if flags.TargetOS != "" {
		existing.TargetOS = flags.TargetOS
	}
	if flags.Instructions != "" {
		existing.Instructions = flags.Instructions
	}
	if flags.Version != "" {
		existing.Version = flags.Version
	}

	// Archive old version before updating
	versionRecord := &model.InstallerVersion{
		ScriptID:     existing.ID,
		Version:      existing.Version,
		TargetOS:     existing.TargetOS,
		Instructions: existing.Instructions,
	}
	if errSave := db.SaveVersion(versionRecord); errSave != nil {
		appErr := apperror.Wrap(errSave, "executeUpdate", map[string]any{"slug": flags.Slug})
		appErr.Code = "E_INSTALLER_UPDATE_FAILED"
		return appErr
	}

	fmt.Printf("Installer %q updated successfully (version: %s, os: %s).\n",
		existing.Name, existing.Version, existing.TargetOS)
	return nil
}

// runInstallerUpdate coordinates flag parsing, database connection, and update execution.
func runInstallerUpdate(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerUpdate", "E_INSTALLER_NIL_COMMAND", map[string]any{
			"error": "command is nil",
		})
	}

	flags, errParse := parseUpdateFlags(args)
	if errParse != nil {
		return errParse
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerUpdate", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerUpdate", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeInstallerUpdate(cmd.Context(), db, flags)
}
