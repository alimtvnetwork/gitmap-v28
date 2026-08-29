// Package cmd — installer_import.go defines the installer import CLI subcommand.
package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// ImportInstallerFlags encapsulates parsed CLI options for installer imports.
type ImportInstallerFlags struct {
	InputPath string
}

// installerImportCmd represents the 'gitmap installer import' subcommand.
var installerImportCmd = &cobra.Command{
	Use:                "import [path]",
	Short:              "Import installer scripts from a zip archive or json file",
	Long:               "Loads installer records and versions from an exported .zip archive or .json file.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerImport(cmd, args)
	},
}

// InstallerImportCmd is an exported alias for installerImportCmd.
var InstallerImportCmd = installerImportCmd

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerImportCmd)
	}
}

// parseImportFlags parses CLI options and positional args for import.
func parseInstallerImportFlags(args []string) (*ImportInstallerFlags, error) {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	fileFlag := fs.String("file", "", "Target file to import from")
	fs.StringVar(fileFlag, "f", "", "File shorthand")

	flagArgs, positional := separateFlagAndPositionalArgs(args)

	if err := fs.Parse(flagArgs); err != nil {
		appErr := apperror.Wrap(err, "parseImportFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"
		return nil, appErr
	}

	targetPath := *fileFlag
	if targetPath == "" && len(positional) > 0 {
		targetPath = positional[0]
	}
	if targetPath == "" {
		targetPath = "gitmap-export.zip"
	}

	return &ImportInstallerFlags{
		InputPath: strings.TrimSpace(targetPath),
	}, nil
}

// importSingleJSON parses and saves a single installer script json payload.
func importSingleJSON(db *store.DB, r io.Reader) error {
	var script model.InstallerScript
	if err := json.NewDecoder(r).Decode(&script); err != nil {
		return err
	}
	if script.Slug == "" {
		script.Slug = slugify(script.Name)
	}
	existing, _ := db.GetInstallerBySlug(script.Slug)
	if existing != nil {
		return nil
	}
	return db.CreateInstaller(&script)
}

// executeImport reads the target archive or json and updates the database.

func executeInstallerImport(ctx context.Context, db *store.DB, flags *ImportInstallerFlags) error {
	if db == nil {
		return apperror.New("executeImport", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "db cannot be nil"})
	}

	if !filepath.IsAbs(flags.InputPath) {
		flags.InputPath, _ = filepath.Abs(flags.InputPath)
	}

	if _, err := os.Stat(flags.InputPath); err != nil {
		appErr := apperror.Wrap(err, "executeImport", map[string]any{"path": flags.InputPath})
		appErr.Code = "E_INSTALLER_FILE_NOT_FOUND"
		return appErr
	}

	if strings.HasSuffix(strings.ToLower(flags.InputPath), ".json") {
		return importFromJSONFile(db, flags.InputPath)
	}

	return importFromZipArchive(flags.InputPath, db)
}

func importFromJSONFile(db *store.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return importSingleJSON(db, f)
}

func importFromZipArchive(path string, db *store.DB) error {
	zr, errZip := zip.OpenReader(path)
	if errZip != nil {
		appErr := apperror.Wrap(errZip, "executeImport", map[string]any{"path": path})
		appErr.Code = "E_INSTALLER_IMPORT_FAILED"
		return appErr
	}
	defer zr.Close()

	count := 0
	for _, file := range zr.File {
		if importZipEntry(db, file) {
			count++
		}
	}

	fmt.Printf("Successfully imported %d installer script(s) from %s.\n", count, path)
	return nil
}

func importZipEntry(db *store.DB, file *zip.File) bool {
	if !strings.HasSuffix(strings.ToLower(file.Name), ".json") {
		return false
	}
	rc, err := file.Open()
	if err != nil {
		return false
	}
	defer rc.Close()
	return importSingleJSON(db, rc) == nil
}

// runInstallerImport coordinates flag parsing, database connection, and import execution.
func runInstallerImport(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return apperror.New("runInstallerImport", "E_INSTALLER_NIL_COMMAND", map[string]any{"error": "command is nil"})
	}

	flags, errParse := parseInstallerImportFlags(args)
	if errParse != nil {
		return errParse
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerImport", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerImport", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeInstallerImport(cmd.Context(), db, flags)
}
