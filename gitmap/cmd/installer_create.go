// Package cmd — installer_create.go defines the installer create CLI subcommand.
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

// CreateInstallerFlags encapsulates parsed CLI options for installer creation.
type CreateInstallerFlags struct {
	Name         string
	Slug         string
	Description  string
	TargetOS     string
	Version      string
	Instructions string
}

// installerCreateCmd represents the 'gitmap installer create' subcommand.
var installerCreateCmd = &cobra.Command{
	Use:                "create <name>",
	Short:              "Create a new installer script record",
	Long:               "Creates and registers a new installer script with target OS and versioning metadata.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerCreate(cmd, args)
	},
}

// InstallerCreateCmd is an exported alias for installerCreateCmd.
var InstallerCreateCmd = installerCreateCmd

// ParseCreateFlags is an exported alias for parseCreateFlags.
var ParseCreateFlags = parseCreateFlags

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerCreateCmd)
	}
}

// slugify converts a human-readable name into a URL/CLI safe slug.
func slugify(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		} else if r == ' ' || r == '_' || r == '.' {
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

// parseCreateFlags parses command-line flags and positional args for installer creation.
func parseCreateFlags(args []string) (*CreateInstallerFlags, error) {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	desc := fs.String("description", "", "Description for the installer script")
	fs.StringVar(desc, "desc", "", "Description shorthand")
	fs.StringVar(desc, "d", "", "Description shorthand")

	targetOS := fs.String("os", "win", "Target operating system (win, ubuntu, centos, debian, arch-linux, all)")
	fs.StringVar(targetOS, "target-os", "win", "Target operating system")
	fs.StringVar(targetOS, "o", "win", "Target operating system shorthand")

	ver := fs.String("version", "v1.0.0", "Initial version for the installer script")
	fs.StringVar(ver, "v", "v1.0.0", "Version shorthand")

	instr := fs.String("instructions", "", "Installation instructions or JSON payload")
	fs.StringVar(instr, "script", "", "Instructions shorthand")
	fs.StringVar(instr, "i", "", "Instructions shorthand")

	slug := fs.String("slug", "", "Custom slug identifier")
	fs.StringVar(slug, "s", "", "Slug shorthand")

	nameFlag := fs.String("name", "", "Name of the installer script")
	fs.StringVar(nameFlag, "n", "", "Name shorthand")

	if err := fs.Parse(args); err != nil {
		appErr := apperror.Wrap(err, "parseCreateFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"
		return nil, appErr
	}

	name := *nameFlag
	if name == "" && fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	if strings.TrimSpace(name) == "" {
		return nil, apperror.New("parseCreateFlags", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "installer name is required",
		})
	}

	finalSlug := *slug
	if finalSlug == "" {
		finalSlug = slugify(name)
	}

	return &CreateInstallerFlags{
		Name:         strings.TrimSpace(name),
		Slug:         finalSlug,
		Description:  strings.TrimSpace(*desc),
		TargetOS:     strings.TrimSpace(*targetOS),
		Version:      strings.TrimSpace(*ver),
		Instructions: strings.TrimSpace(*instr),
	}, nil
}

// executeCreate executes the persistence logic for creating an installer in the DB.
//nolint:revive
func executeCreate(ctx context.Context, db *store.DB, flags *CreateInstallerFlags) error {
	if db == nil {
		return apperror.New("executeCreate", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "db cannot be nil",
		})
	}
	if flags == nil {
		return apperror.New("executeCreate", "E_INSTALLER_INVALID_INPUT", map[string]any{
			"error": "flags cannot be nil",
		})
	}

	script := &model.InstallerScript{
		Name:         flags.Name,
		Slug:         flags.Slug,
		Description:  flags.Description,
		TargetOS:     flags.TargetOS,
		Version:      flags.Version,
		Instructions: flags.Instructions,
	}

	if err := db.CreateInstaller(script); err != nil {
		appErr := apperror.Wrap(err, "executeCreate", map[string]any{
			"name": flags.Name,
			"slug": flags.Slug,
		})
		appErr.Code = "E_INSTALLER_CREATE_FAILED"
		return appErr
	}

	fmt.Printf("Installer %q created successfully (slug: %s, version: %s, os: %s).\n",
		script.Name, script.Slug, script.Version, script.TargetOS)
	return nil
}

// runInstallerCreate coordinates flag parsing, database connection, and creation execution.
func runInstallerCreate(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerCreate", "E_INSTALLER_NIL_COMMAND", map[string]any{
			"error": "command is nil",
		})
	}

	flags, errParse := parseCreateFlags(args)
	if errParse != nil {
		return errParse
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerCreate", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerCreate", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeCreate(cmd.Context(), db, flags)
}
