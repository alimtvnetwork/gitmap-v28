// Package cmd — installer_update_win.go defines the installer update-win CLI subcommand.
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

// UpdateWinInstallerFlags encapsulates parsed CLI options for windows installer updates.
type UpdateWinInstallerFlags struct {
	Slug         string
	Description  string
	Version      string
	Instructions string
}

// installerUpdateWinCmd represents the 'gitmap installer update-win' subcommand.
var installerUpdateWinCmd = &cobra.Command{
	Use:                "update-win <slug>",
	Short:              "Update an existing installer script record specifically for Windows",
	Long:               "Updates an installer script record for Windows target OS and records version history.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerUpdateWin(cmd, args)
	},
}

// InstallerUpdateWinCmd is an exported alias for installerUpdateWinCmd.
var InstallerUpdateWinCmd = installerUpdateWinCmd

// ParseUpdateWinFlags is an exported alias for parseUpdateWinFlags.
var ParseUpdateWinFlags = parseUpdateWinFlags

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerUpdateWinCmd)
	}
}

// parseUpdateWinFlags parses command-line flags and positional args for windows installer update.
func parseUpdateWinFlags(args []string) (*UpdateWinInstallerFlags, error) {
	fs := flag.NewFlagSet("update-win", flag.ContinueOnError)
	desc := fs.String("description", "", "Description for the installer script")
	fs.StringVar(desc, "desc", "", "Description shorthand")
	fs.StringVar(desc, "d", "", "Description shorthand")

	ver := fs.String("version", "", "Version for the installer script")
	fs.StringVar(ver, "v", "", "Version shorthand")

	instr := fs.String("instructions", "", "Installation instructions or JSON payload")
	fs.StringVar(instr, "script", "", "Instructions shorthand")
	fs.StringVar(instr, "i", "", "Instructions shorthand")

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
		appErr := apperror.Wrap(err, "parseUpdateWinFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"
		return nil, appErr
	}

	slug := *slugFlag
	if slug == "" && len(positional) > 0 {
		slug = positional[0]
	}

	if strings.TrimSpace(slug) == "" {
		return nil, apperror.New("parseUpdateWinFlags", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "installer slug is required",
		})
	}

	return &UpdateWinInstallerFlags{
		Slug:         strings.TrimSpace(slug),
		Description:  strings.TrimSpace(*desc),
		Version:      strings.TrimSpace(*ver),
		Instructions: strings.TrimSpace(*instr),
	}, nil
}

// executeUpdateWin executes the update logic specifically for Windows target OS.
func executeUpdateWin(ctx context.Context, db *store.DB, flags *UpdateWinInstallerFlags) error {
	if db == nil {
		return apperror.New("executeUpdateWin", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "db cannot be nil",
		})
	}
	if flags == nil {
		return apperror.New("executeUpdateWin", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "flags cannot be nil",
		})
	}

	existing, errGet := db.GetInstallerBySlug(flags.Slug)
	if errGet != nil {
		appErr := apperror.Wrap(errGet, "executeUpdateWin", map[string]any{"slug": flags.Slug})
		appErr.Code = "E_INSTALLER_NOT_FOUND"
		return appErr
	}

	if flags.Description != "" {
		existing.Description = flags.Description
	}
	existing.TargetOS = "win"
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
		TargetOS:     "win",
		Instructions: existing.Instructions,
	}
	if errSave := db.SaveVersion(versionRecord); errSave != nil {
		appErr := apperror.Wrap(errSave, "executeUpdateWin", map[string]any{"slug": flags.Slug})
		appErr.Code = "E_INSTALLER_UPDATE_FAILED"
		return appErr
	}

	fmt.Printf("Windows Installer %q updated successfully (version: %s).\n",
		existing.Name, existing.Version)
	return nil
}

// runInstallerUpdateWin coordinates flag parsing, database connection, and update execution.
func runInstallerUpdateWin(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerUpdateWin", "E_INSTALLER_NIL_COMMAND", map[string]any{
			"error": "command is nil",
		})
	}

	flags, errParse := parseUpdateWinFlags(args)
	if errParse != nil {
		return errParse
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerUpdateWin", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerUpdateWin", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeUpdateWin(cmd.Context(), db, flags)
}
